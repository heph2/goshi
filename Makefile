PREFIX =	/usr/local/
DESTDIR =

.PHONY: all goshi clean

all: goshi

goshi:
	cd cmd/goshi && go build

clean:
	rm -f goshi

install: goshi
	mkdir -p ${DESTDIR}${PREFIX}/bin/
	cd cmd/goshi &&	install -m 0555 goshi ${DESTDIR}${PREFIX}/bin
