export namespace backend {
	
	export class BookStats {
	    probablyKnownWords: number;
	    knownCharacters: number;
	    totalCharacters: number;
	    totalWords: number;
	    totalKnownWords: number;
	    targets: number[];
	    targetOccurances: number[];
	    needToKnow: number[];
	
	    static createFrom(source: any = {}) {
	        return new BookStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.probablyKnownWords = source["probablyKnownWords"];
	        this.knownCharacters = source["knownCharacters"];
	        this.totalCharacters = source["totalCharacters"];
	        this.totalWords = source["totalWords"];
	        this.totalKnownWords = source["totalKnownWords"];
	        this.targets = source["targets"];
	        this.targetOccurances = source["targetOccurances"];
	        this.needToKnow = source["needToKnow"];
	    }
	}
	export class Book {
	    author: string;
	    title: string;
	    cover: string;
	    filepath: string;
	    bookId: number;
	    favorite: boolean;
	    hasRead: boolean;
	    stats: BookStats;
	
	    static createFrom(source: any = {}) {
	        return new Book(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.author = source["author"];
	        this.title = source["title"];
	        this.cover = source["cover"];
	        this.filepath = source["filepath"];
	        this.bookId = source["bookId"];
	        this.favorite = source["favorite"];
	        this.hasRead = source["hasRead"];
	        this.stats = this.convertValues(source["stats"], BookStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class DictionaryDefinition {
	    definition: string;
	    pronunciation: string;
	
	    static createFrom(source: any = {}) {
	        return new DictionaryDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.definition = source["definition"];
	        this.pronunciation = source["pronunciation"];
	    }
	}
	export class Fields {
	    word: string;
	    sentence: string;
	    englishDefn: string;
	    chineseDefn: string;
	    pinyin: string;
	    imageUrls: string[];
	
	    static createFrom(source: any = {}) {
	        return new Fields(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.sentence = source["sentence"];
	        this.englishDefn = source["englishDefn"];
	        this.chineseDefn = source["chineseDefn"];
	        this.pinyin = source["pinyin"];
	        this.imageUrls = source["imageUrls"];
	    }
	}
	export class FlaggedCard {
	    word: string;
	    sentence: string;
	
	    static createFrom(source: any = {}) {
	        return new FlaggedCard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.sentence = source["sentence"];
	    }
	}
	export class ImageInfo {
	    thumbnailUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.thumbnailUrl = source["thumbnailUrl"];
	    }
	}
	export class Problems {
	    Flagged: boolean;
	    MissingImage: boolean;
	    MissingSentence: boolean;
	    MissingSentenceAudio: boolean;
	    MissingWordAudio: boolean;
	    MissingPinyin: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Problems(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Flagged = source["Flagged"];
	        this.MissingImage = source["MissingImage"];
	        this.MissingSentence = source["MissingSentence"];
	        this.MissingSentenceAudio = source["MissingSentenceAudio"];
	        this.MissingWordAudio = source["MissingWordAudio"];
	        this.MissingPinyin = source["MissingPinyin"];
	    }
	}
	export class ProblemCard {
	    Word: string;
	    Problems: Problems;
	
	    static createFrom(source: any = {}) {
	        return new ProblemCard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Word = source["Word"];
	        this.Problems = this.convertValues(source["Problems"], Problems);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class RawAnkiNote {
	    noteId: number;
	    fields: Fields;
	
	    static createFrom(source: any = {}) {
	        return new RawAnkiNote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.noteId = source["noteId"];
	        this.fields = this.convertValues(source["fields"], Fields);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UnknownWordEntry {
	    word: string;
	    occurance?: number;
	    definition?: string;
	    pinyin?: string;
	
	    static createFrom(source: any = {}) {
	        return new UnknownWordEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.occurance = source["occurance"];
	        this.definition = source["definition"];
	        this.pinyin = source["pinyin"];
	    }
	}
	
	export class WordOccuranceRow {
	    word: string;
	    occurance: number;
	
	    static createFrom(source: any = {}) {
	        return new WordOccuranceRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.occurance = source["occurance"];
	    }
	}
	export class WordStats {
	    words: number;
	    characters: number;
	
	    static createFrom(source: any = {}) {
	        return new WordStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.words = source["words"];
	        this.characters = source["characters"];
	    }
	}

}

