# Planned and Completed Features

Think of the todos as just ideas I would like to optimistically work to, and not actual roadmapped features

### Todo
- [ ] 90% test coverage for go code
- [ ] unflag anki cards on edit
- [ ] allow mapping output to arbitrary anki note format
- [ ] free edit of card fields
- [ ] card maintenance (search for missing image, audio, etc)
- [ ] unified pinyin generation
- [ ] use forvo for term audio [unclear if possible]
- [ ] get example sentences from some online service [if free ones exist] 
- [ ] support frequency lists
- [ ] chinese voice picker
- [ ] linux scrollbars dont work in webview
- [ ] track book progress (eg 15% read)
- [ ] update image / audio when editing card
- [ ] write dev instal guide 
- [ ] use calibre (if possible) to convert from epub etc to txt file
- [ ] store our own copy of image cover, book text, etc
- [ ] more compact library table view
- [ ] manage known words list


### In Progress

These will definitely be implemented and have priority

- [ ] allow cc-cedict to be installed on first start
- [ ] interactive workflow for importing from calibre (control which books to import etc)
- [ ] import user word list from various sources
- [ ] add new book (without needing calibre)


### DONE
- [x] Move from naive-ui to daisyui
- [x] Rewrite using golang / wails.io
- [x] Segment text only using words in user dictionary
- [x] Speed up storage and loading of segmented text
- [x] Generate audio for sentences using Azure
- [x] Load images for cards from web
- [x] Manage api-keys for image and tts services
- [x] Manage dictionaries
- [x] Start rewrite as electron application
