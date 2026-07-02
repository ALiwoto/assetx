# Converters
All format converters code will go here.

## tgsConverter

The TGS converter handles Telegram WebM/VP9 emoji and sticker files and writes a PNG sprite sheet.

Sprite sheets are packed in row-major order: frames go left-to-right across the first row, then continue left-to-right on the next row.

Example with 10 frames and 4 columns:

```text
1,  2,  3,  4
5,  6,  7,  8
9, 10
```



