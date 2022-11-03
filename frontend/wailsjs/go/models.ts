export namespace backend {
	
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

