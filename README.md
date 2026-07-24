# assetx

Generate game assets with AI agents without hitting your head to the wall.

## Example command

```bash
assetx image --model "gpt-image-2" --background "transparent" --prompt "create a battle win header" --avoid "watermark" --avoid "text artifacts" --example "example1.png" --example-note "screenshot of my game" --example "example2.png" --example-note "existing UI asset style reference" --quality medium --size 1024x1024 --out assets/sprites/slime.png
```

Convert a Telegram WebM/VP9 emoji or sticker file into a PNG sprite sheet:

```bash
assetx convert-tgs --in "input_file.tgs" --out "output_file.png"
```

Convert a WebP image into PNG:

```bash
assetx convert-webp --in "input_file.webp" --out "output_file.png"
assetx convert-webp "input_file.webp"
```

Search the web through the OpenAI Responses API and print sourced Markdown:

```bash
assetx search --domain fab.com --query "Find modular medieval character systems for Unreal Engine"
```

