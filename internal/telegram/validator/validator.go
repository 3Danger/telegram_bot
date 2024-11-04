package validator

import (
	"time"

	tele "github.com/PaulSonOfLars/gotgbot/v2"
)

type Bound[T int | int64 | time.Duration] struct {
	Max, Min T
}

func NewBound[T int | int64 | time.Duration](min, max T) Bound[T] {
	return Bound[T]{
		Max: max,
		Min: min,
	}
}

type MediaValidator struct {
	photoResolution Bound[int64]
	photoFileSize   Bound[int64]
	videoResolution Bound[int64]
	videoFileSize   Bound[int64]
	videoDuration   Bound[time.Duration]
}

func New(
	photoResolution Bound[int64],
	photoFileSize Bound[int64],
	videoResolution Bound[int64],
	videoFileSize Bound[int64],
	videoDuration Bound[time.Duration],
) *MediaValidator {
	return &MediaValidator{
		photoResolution: photoResolution,
		photoFileSize:   photoFileSize,
		videoResolution: videoResolution,
		videoFileSize:   videoFileSize,
		videoDuration:   videoDuration,
	}
}

type Err struct {
	msg    string
	tooBig bool
}

func (e *Err) Error() string { return e.msg }
func (e *Err) TooBig() bool  { return e.tooBig }

func newErr(msg string, tooBig bool) error {
	return &Err{
		msg:    msg,
		tooBig: tooBig,
	}
}

var (
	ErrPhotoIsTooBigInHeight   = newErr("the photo is too big in height", true)
	ErrPhotoIsTooBigInWidth    = newErr("the photo is too big in width", true)
	ErrPhotoFileSizeIsTooBig   = newErr("the photo file size is too big", true)
	ErrPhotoIsTooSmallInHeight = newErr("the photo is too small in height", false)
	ErrPhotoIsTooSmallInWidth  = newErr("the photo is too small in width", false)
	ErrPhotoFileSizeIsTooSmall = newErr("the photo file size is too small", false)
)

func (v *MediaValidator) ValidatePhoto(f *tele.PhotoSize) error {
	switch {
	case v.photoResolution.Max < f.Height:
		return ErrPhotoIsTooBigInHeight
	case v.photoResolution.Max < f.Width:
		return ErrPhotoIsTooBigInWidth
	case v.photoFileSize.Max < f.FileSize:
		return ErrPhotoFileSizeIsTooBig
	}

	switch {
	case v.photoResolution.Min > f.Height:
		return ErrPhotoIsTooSmallInHeight
	case v.photoResolution.Min > f.Width:
		return ErrPhotoIsTooSmallInWidth
	case v.photoFileSize.Min > f.FileSize:
		return ErrPhotoFileSizeIsTooSmall
	}

	return nil
}

var (
	ErrVideoIsTooBigInHeight       = newErr("the video is too big in height", true)
	ErrVideoIsTooBigInWidth        = newErr("the video is too big in width", true)
	ErrVideoFileSizeIsTooBig       = newErr("the video file size is too big", true)
	ErrVideoFileDurationIsTooLong  = newErr("the video duration is too long", true)
	ErrVideoIsTooSmallInHeight     = newErr("the video is too small in height", false)
	ErrVideoIsTooSmallInWidth      = newErr("the video is too small in width", false)
	ErrVideoFileSizeIsTooSmall     = newErr("the video file size is too small", false)
	ErrVideoFileDurationIsTooShort = newErr("the video duration is too short", true)
)

func (v *MediaValidator) ValidateVideo(f *tele.Video) error {
	switch {
	case v.videoResolution.Max < f.Height:
		return ErrVideoIsTooBigInHeight
	case v.videoResolution.Max < f.Width:
		return ErrVideoIsTooBigInWidth
	case v.videoFileSize.Max < f.FileSize:
		return ErrVideoFileSizeIsTooBig
	case int64(v.videoDuration.Max.Seconds()) < f.Duration:
		return ErrVideoFileDurationIsTooLong
	}

	switch {
	case v.videoResolution.Min > f.Height:
		return ErrVideoIsTooSmallInHeight
	case v.videoResolution.Min > f.Width:
		return ErrVideoIsTooSmallInWidth
	case v.videoFileSize.Min > f.FileSize:
		return ErrVideoFileSizeIsTooSmall
	case int64(v.videoDuration.Min.Seconds()) > f.Duration:
		return ErrVideoFileDurationIsTooShort
	}

	return nil
}

func (v *MediaValidator) ValidateVideoNote(f *tele.VideoNote) error {
	switch {
	case v.videoFileSize.Max < f.FileSize:
		return ErrVideoFileSizeIsTooBig
	case int64(v.videoDuration.Max.Seconds()) < f.Duration:
		return ErrVideoFileDurationIsTooLong
	}

	switch {
	case v.videoFileSize.Min > f.FileSize:
		return ErrVideoFileSizeIsTooSmall
	case int64(v.videoDuration.Min.Seconds()) > f.Duration:
		return ErrVideoFileDurationIsTooShort
	}

	return nil
}