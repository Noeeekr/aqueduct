package server

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Noeeekr/aqueduct/internal"
	"github.com/gin-gonic/gin"
)

type AnchorInfo struct {
	Path     string
	FileName string
	Session  string // If created in session gets the path of the file to render del btn
}

type Handlers struct {
	*internal.Info
	Logger *internal.Logger
}

func NewHandlers() *Handlers {
	return &Handlers{
		Info:   internal.NewInstance().Info,
		Logger: internal.NewInstance().Logger,
	}
}

// SERVE TEMPLATE FILE
func (h *Handlers) ServeTemplate(ctx *gin.Context) {
	_path, err := getRelativePath(ctx)
	if err != nil {
		ctx.Redirect(301, "/?path=/")
		return
	}

	// Get absolute path to shared folder
	files_root := h.SharedFolder
	folder := path.Join(files_root, _path)

	// Handle the case where the shared folder doesnt exist
	folder_content, err := os.Stat(folder)
	if err != nil {
		ctx.Redirect(301, "/?path=/")
		return
	}

	// Handle the case where it is a file
	if !folder_content.IsDir() {
		ctx.Redirect(301, path.Join("/files", _path))
		return
	}

	folders, err := os.ReadDir(folder)
	if err != nil {
		ctx.Redirect(302, "/notfound")
		return
	}

	// Handle getting cookie to check if file is part of session
	cookie := parseSessionCookie(ctx)

	// Read dir|files names to use in rendering folder navigation
	var folder_names []AnchorInfo

	for _, dir := range folders {
		del_btn_path := ""
		dir_path := path.Join(_path, dir.Name())

		if (*cookie)[path.Join(files_root, dir_path)] {
			del_btn_path = dir_path
		}

		folder_names = append(folder_names, AnchorInfo{
			FileName: dir.Name(),
			Path:     dir_path,
			Session:  del_btn_path,
		})
	}

	// Handle the case where there's no parent folder to go
	var last_path string = "/"
	if len(_path) > 1 {
		last_path = path.Join(_path, "..")
	}

	data := gin.H{
		"path":      _path,
		"last_path": last_path,
		"folders":   folder_names,
	}

	ctx.HTML(http.StatusOK, "index.html", data)
}

func (h *Handlers) HandleUpload(ctx *gin.Context) {
	ctx.Request.ParseMultipartForm(100 << 20)
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
	}

	// Relative path in folder structure
	_path, err := getRelativePath(ctx)
	if err != nil {
		ctx.Redirect(301, "/?path=/")
	}

	root_folder := h.SharedFolder

	// Holds all the dirs and files created
	var cookie_val *map[string]bool = &map[string]bool{}

	for _, files := range form.File["files"] {
		// If cannot get file name returns err
		if files.Header.Get("Content-Disposition") == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		} else {
			// Get header fields [key|val]
			fields := strings.Split(files.Header.Get("Content-Disposition"), ";")

			for _, raw_field := range fields {
				// Separes header fields [key|val] into an array
				field := strings.Split(raw_field, "=")
				// Get only valid [key|val] fields and find the filename header
				if len(field) == 2 && strings.Contains(field[0], "filename") {
					file_path := path.Join(_path, strings.Trim(field[1], `"`))

					// Check if filepath is .. or any kind of troll/../../../..
					if !isSecurePath(file_path) {
						ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "O arquivo " + file_path + " é um risco de segurança."})
						return
					}

					// Check if file already exists
					if file_info, err := os.Stat(path.Join(root_folder, file_path)); err == nil {
						if !file_info.IsDir() {
							ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao criar o arquivo " + files.Filename + ". O arquivo já existe e o administrador optou por não permitir adulterações."})
							return
						}
					}

					// Create dirs if necessary
					if file_path != files.Filename {
						err := os.MkdirAll(path.Join(root_folder, file_path, ".."), 0666)

						paths := path.Join(root_folder, file_path)

						for strings.HasPrefix(paths, path.Join(root_folder, _path)) && paths != path.Join(root_folder, _path) {
							(*cookie_val)[paths] = true

							paths = path.Join(paths, "..")
						}

						if err != nil {
							ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao criar o caminho de pasta " + path.Join(root_folder, file_path, "..")})
							return
						}

					}

					// Create file and copy content to it
					outFile, err := os.Create(path.Join(root_folder, file_path))
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao criar o arquivo" + files.Filename})
						return
					}

					inFile, err := files.Open()
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao ler o arquivo" + files.Filename})
						return
					}

					if _, err = io.Copy(outFile, inFile); err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao copiar o arquivo" + files.Filename + "para o servidor"})
						return
					}

					inFile.Close()
					outFile.Close()
				}
			}
		}
	}

	var cookie *map[string]bool = &map[string]bool{}

	c, err := ctx.Cookie("session")
	if err != nil {
		cookie = &map[string]bool{}
	} else {
		err = json.Unmarshal([]byte(c), cookie)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao salvar o histórico da sessão [208]"})
			return
		}
	}

	for key := range *cookie_val {
		(*cookie)[key] = true
	}

	_cookie, err := json.Marshal(cookie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Falha ao salvar o histórico da sessão [209]"})
		return
	}

	ctx.SetCookie(
		"session",
		string(_cookie),
		3600,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": ""})
}

func (h *Handlers) HandleDelete(ctx *gin.Context) {
	relative_path, err := getRelativePath(ctx)
	if err != nil {
		ctx.JSON(http.StatusCreated, gin.H{"success": false, "message": "Caminho relativo não especificado"})
		return
	}

	root_path := h.Info.SharedFolder
	absolute_path := path.Join(root_path, relative_path)

	if err := os.RemoveAll(absolute_path); err != nil {
		ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Falha ao deletar a pasta. Error:" + err.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": ""})
}

// Get the URL path, which is the path starting that goes after the root_folder path
func getRelativePath(ctx *gin.Context) (path string, err error) {
	_path := ctx.Request.URL.Query().Get("path")

	if !isSecurePath(_path) {
		return "/", errors.New(" Caminho de arquivo inválido")
	}

	return _path, nil
}

// Check if the path contains "../" or "/.."
func isSecurePath(_path string) bool {
	if _path == "" {
		return false
	}

	if len(_path) < 1 && _path != "/" {
		return false
	}

	if strings.Contains(_path, "../") || strings.HasPrefix(_path, "..") || strings.Contains(_path, "/..") {
		return false
	}

	return true
}

// Revalidates the cookie max age so it doesn't expire
func (h *Handlers) HandleCookie(ctx *gin.Context) {
	cookie, err := ctx.Cookie("session")
	if err != nil {
		c, err := json.Marshal(map[string]bool{})
		if err != nil {
			c = []byte("")
		}
		ctx.SetCookie(
			"session",
			string(c),
			3600,
			"/",
			"",
			false,
			true,
		)
	} else {
		ctx.SetCookie(
			"session",
			cookie,
			3600,
			"/",
			"",
			false,
			true,
		)
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": ""})
}

func parseSessionCookie(ctx *gin.Context) *map[string]bool {
	var cookie *map[string]bool = &map[string]bool{}

	// Get cookie
	c, err := ctx.Cookie("session")
	if err != nil {
		// Forge a new empty cookie
		_cookie, err := json.Marshal(map[string]string{})
		if err != nil {
			_cookie = []byte("{}")
		}

		c = string(_cookie)
	}

	// Parse cookie
	err = json.Unmarshal([]byte(c), cookie)
	if err != nil {
		cookie = &map[string]bool{}
	}

	return cookie
}

func (h *Handlers) HandleDownload(ctx *gin.Context) {
	relative_path, err := getRelativePath(ctx)
	if err != nil {
		ctx.JSON(http.StatusCreated, gin.H{"success": false, "message": "Caminho relativo não especificado"})
		return
	}

	root_path := h.Info.SharedFolder
	relative_path, folder_name := path.Split(relative_path)
	absolute_path := path.Join(root_path, relative_path)

	zipFolder(
		ctx,
		fmt.Sprintf("%s/%s", absolute_path, folder_name),
	)

	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", folder_name))
}

func zipFolder(ctx *gin.Context, folder_path string) error {
	// Creates a writer to write to buffer
	zipWriter := zip.NewWriter(ctx.Writer)

	// Go through all files in folder
	err := filepath.WalkDir(folder_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		zipFilePath := strings.TrimPrefix(path, folder_path+"/")
		zipFile, err := zipWriter.Create(zipFilePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, file)

		return err
	})

	if err != nil {
		return err
	}

	return zipWriter.Close()
}
