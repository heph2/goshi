# Goshi

Tool for concurrently download Manga Scans.
Actually the following sites are supported:
- Mangaeden (currently Italian only)
- Manganelo
- Mangaworld

## Getting Started

### Installing

    git clone https://git.mrkeebs.eu/goshi
    cd goshi
    make
    sudo make install
    

### Usage

Fetching the available Chapters scraping manganelo

    goshi [-fetch] <manga> [-scraper] manganelo
        
Downloading a single episode using the ID associated to the chapter returned by -fetch

    goshi [-fetch] <manga> [-scraper] manganelo -down <ID>
