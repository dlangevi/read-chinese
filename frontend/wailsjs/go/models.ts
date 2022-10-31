export namespace backend {
	
	export class WordRow {
	    word: string;
	    occurance: number;
	
	    static createFrom(source: any = {}) {
	        return new WordRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.occurance = source["occurance"];
	    }
	}

}

