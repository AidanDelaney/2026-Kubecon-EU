all: serve

presentation:
	pandoc -t revealjs -s presentation.md -o presentation.html \
		-V theme=moon -V transition=none -V controls=true -V progress=true -V slideNumber=true \
		-V revealjs-url=https://unpkg.com/reveal.js@^4 -V navigationMode=linear \
		--css tweaks.css \
		--slide-level=2

serve: presentation
	python3.14 -m http.server