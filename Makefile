all: serve

presentation:
	pandoc -t revealjs -s presentation.md -o presentation.html \
		-V theme=white -V transition=none -V controls=true -V progress=true -V slideNumber=true \
		-V revealjs-url=https://unpkg.com/reveal.js@^4 -V navigationMode=linear \
		--css tweaks.css \
		--slide-level=2

pdf: presentation
	pandoc presentation.md -t beamer -o presentation.pdf \
		-V geometry:margin=1cm \
		-V fontsize=12pt \
		-V colorlinks=true \
		-V linkcolor=blue \
		-V urlcolor=blue

serve: presentation
	python3.14 -m http.server