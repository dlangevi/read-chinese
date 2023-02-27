# read-chinese
## Blocking issues
Currently the following have to be implemented before I would recommend any
random person on the internet tries to use this

* I need to make a guide on how to get your own Azure api keys
* I need to make a guide on how to import known word lists

## What is this?
### Who is this for
Intermediate to Advanced learners with 
* a vocab of (at least) 1000 characters, and 2000 known words.
* an interest in reading 
* experience using Anki

### How it works
For a brief overview / demo of the project (v0.0.3) you can watch 
[this video](https://youtu.be/dAqE2YquGlw "read-chinese demo") the first 5
minutes contains most of the useful info.

### Technical Overview
A (technical) summary of the workflow is as follows

* Import txt files containing the text of a book
* Run Jieba segmentation algorithm on text (filter valid words it can segment
  to those in user dictionary)
* Using knowledge of users known words, create a set of all T1 sentences in all books
* Present the user list of most frequent unknown words in the text
* For each word, provide 
    * List of T1 sentences containing this word
    * Definitions from dictionaries
    * Images from Microsoft Azure image search
    * Generated azure TTS of the word / sentence (if desired)

### Getting Started

#### Setup (I hope to automate much of this soon)
1. Acquire a ebook copies of the books you want to read. (If you legally own a
 physical copy, look on anna's archive for copy you can use)
2. Load these books into calibre and 'convert' them to txt format.
3. Sync the books from calibre to read-chinese
4. Import a dictionary (cc-cedict can be imported by just clicking a button)
5. (Ideally) Import a list of new words
6. Generate your own azure api keys for the free image and tts services
7. Open anki (with anki-connect plugin installed)

#### Making the flashcards

1. Click on the book you want to study for
2. Click 'Make Flashcards' button
3. Go go go

### What is this trying to accomplish
I have spent a much of my time learning Chinese through reading, as reading is
a hobby I enjoy outside of language learning, and I think it is a great way to
pick up new vocab. I don't like reading using a computer, and I tend to like
smaller phones, so I have tended to stick to reading on either a kindle or
physical books. This makes it harder to lookup the pronunciation / meaning of
new words, so I rely a bit more heavily on 'prestudying' for a book, making
sure to learn the highest frequency unknown words and characters ahead of time.

Early on the biggest time sink for this was actually creating the flashcards.
Even using other available tools for this I would spend up to 2+ hours creating
40 high quality flashcards. Since I would also try to review 20 new cards a day
at the time, this meant I was spending at least an hour each day just making
cards. This felt like a pretty big waste of time so I wanted to focus on making
a program that would speed up this process for me.

At this point it feels like I can make those same 40 flashcards in only 20
minutes now, so I think it does what I need pretty well.

### Some Drawbacks
I think the biggest criticism I would have of this is approach is that the long
term benefit of creating flashcards for only the most frequent words is that
you would probably acquire these words anyways if you just read the book. 

There might be a more efficient order in which prioritize words (and would be
happy to implement such an ordering) but for now more frequent is easiest, since
those words will tend to have more T1 sentences to choose from.

## Development
### Dependancies:
Let me know if something else is needed that I am not aware of

* [calibre](https://calibre-ebook.com/download)
* [wails](https://wails.io/docs/gettingstarted/installation) (make sure to add $HOME/go/bin to path)
* [yarn](https://yarnpkg.com/getting-started/install)
* some sort of c++ compiler (to compile the Jieba implementation)

### to run dev mode:
`wails dev`

### to build
`wails build`

Currently I am using the following flags to create the windows installer. This
is because I primarily develop and test on linux and want to get the most info
possible out of problems when they happen on windows

`wails build -nsis -windowsconsole -debug`
