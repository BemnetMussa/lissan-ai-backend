package handler

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"lissanai.com/backend/internal/domain/interfaces"
)

type PronunciationHandler struct {
	pronunciationUC interfaces.PronunciationUsecase
}

func NewPronunciationHandler(uc interfaces.PronunciationUsecase) *PronunciationHandler {
	return &PronunciationHandler{pronunciationUC: uc}
}

// GetSentences handles the REST request to fetch the list of practice sentences.
func (h *PronunciationHandler) GetSentences(c *gin.Context) {
	sentences := h.pronunciationUC.GetPracticeSentences()
	c.JSON(http.StatusOK, sentences)
}

// AssessPronunciation handles the file upload and text data for assessment.
func (h *PronunciationHandler) AssessPronunciation(c *gin.Context) {
	targetText := c.PostForm("target_text")
	if targetText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Form field 'target_text' is required"})
		return
	}

	audioFile, err := c.FormFile("audio_data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file 'audio_data' is required"})
		return
	}

	openedFile, err := audioFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open audio file"})
		return
	}
	defer openedFile.Close()

	audioData, err := ioutil.ReadAll(openedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read audio file"})
		return
	}

	log.Printf("Received request to assess pronunciation for text: '%s' with audio size %d bytes.", targetText, len(audioData))

	feedback, err := h.pronunciationUC.AssessPronunciation(c.Request.Context(), targetText, audioData)
	if err != nil {
		log.Printf("Error from pronunciation usecase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assess pronunciation"})
		return
	}

	c.JSON(http.StatusOK, feedback)
}
