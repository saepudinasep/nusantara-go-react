package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/pkg/response"
)

type TagihanHandler struct {
	tagihanUsecase domain.TagihanUsecase
}

func NewTagihanHandler(tagihanUsecase domain.TagihanUsecase) *TagihanHandler {
	return &TagihanHandler{tagihanUsecase: tagihanUsecase}
}

type tagihanResponse struct {
	Siswa   *domain.Student     `json:"siswa"`
	Tagihan []domain.SppTagihan `json:"tagihan"`
}

// GetTagihan menangani GET /api/siswa/tagihan — status Lunas/Belum per bulan untuk SEMUA jenis SPP
func (h *TagihanHandler) GetTagihan(c *gin.Context) {
	student, tagihan, err := h.tagihanUsecase.GetTagihan(c.Request.Context(), currentUserID(c))
	if err != nil {
		if errors.Is(err, domain.ErrStudentProfileMissing) {
			response.Error(c, http.StatusNotFound, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "gagal mengambil data tagihan")
		return
	}

	response.Success(c, http.StatusOK, "berhasil mengambil data tagihan", tagihanResponse{
		Siswa:   student,
		Tagihan: tagihan,
	})
}
