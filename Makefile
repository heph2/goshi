PREFIX =	/usr/local/
DESTDIR =

.PHONY: all goshi clean

all: goshi

goshi:
	go build -o cmd/goshi/goshi ./cmd/goshi

clean:
	rm -f goshi

install: goshi
	mkdir -p ${DESTDIR}${PREFIX}/bin/
	cd cmd/goshi &&	install -m 0555 goshi ${DESTDIR}${PREFIX}/bin
