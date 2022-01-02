package handlers

// TODO v2 smarter watermark
// e.g random time , random pos, random  size , random color , match with backgroud color , match transparency

// dependencies
// FFMPEG
// Aubio
// ImageMagick
// 7Zip

//ffmpeg audio watermark

// get audio duriation
// ffprobe -i Audio_Sample.mp3  -show_entries format=duration -v quiet -of csv="p=0"

// get audio BPM
// aubio tempo Audio_Sample.mp3

// Watermark audio
//ffmpeg -i main.mp3 -filter_complex "amovie=beep.wav:loop=0,asetpts=N/SR/TB,adelay=10s:all=1[beep]; [0][beep]amix=duration=shortest,volume=2"   out.mp3
//ffmpeg -i main.mp3 '-filter_complex', '[0:a]volume=volume=1[aout0];[1:a]volume=volume=2[aout1];[aout1]aloop=loop=-1:size=2e+09,adelay=2000,atrim=start=0:end=2:duration=6[aconcat];[aout0][aconcat]amix=inputs=2:duration=longest:dropout_transition=4 [aout]',

// PDF

// No Perm
// pdfengine encrypt -m aes -k 256 -perm none -upw userdecryptkey -opw ownerdecryptkey sample_pdf.pdf sample_encrypted.pdf
// FullPerm
// pdfengine encrypt -m aes -k 256 -perm all -upw userdecryptkey -opw ownerdecryptkey sample_pdf.pdf sample_encrypted.pdf

// Video

// Get the size
// ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 sample_video.mov
// rezie the watermark image
// ffmpeg -i sample_video.mov -i pazari-resized_15.png -filter_complex "overlay=x=(main_w-overlay_w)/2:y=(main_h-overlay_h)/2" output.mp4

// Image
// the the size
// magick identify -ping -format '=> %w %h' sample_image.png
// Resize the watermark to fit
// magick convert pazari-darkest.png -resize 832x720  pazari-resized.png
// Do the watermark
// magick composite  -dissolve 15% -gravity SouthWest  pazari-resized.png sample_image.png sample_image_watermarked.png

//7z a -t7z -m0=lzma2 -mx=9 -mfb=64 -md=32m -ms=on -mhe=on -p"password"  output.zip sample_pdf.pdf
//7z l -slt output.zip
