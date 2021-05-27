PREFIX =	/usr/local/
DESTDIR =

.PHONY: all goshi clean

all: goshi

saturn:
	go build

clean:
	rm -f goshi

install: goshi
	mkdir -p ${DESTDIR}${PREFIX}/bin/
	install -m 0555 goshi ${DESTDIR}${PREFIX}/bin
