package internal

var Images = make(map[int32]ImageData) // negative = crops; 0 = White1x1; positive = full images
var NextImageId int16
var NextImageCropId int16
