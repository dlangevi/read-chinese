export namespace backend {
	
	
	export class DictionaryEntry {
	    definition: string;
	    pronunciation: string;
	
	    static createFrom(source: any = {}) {
	        return new DictionaryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.definition = source["definition"];
	        this.pronunciation = source["pronunciation"];
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

