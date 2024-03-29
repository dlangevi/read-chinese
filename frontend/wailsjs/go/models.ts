export namespace backend {
	
	export class FieldsMapping {
	    firstField: string;
	    hanzi: string;
	    exampleSentence: string;
	    englishDefinition: string;
	    chineseDefinition: string;
	    pinyin: string;
	    hanziAudio: string;
	    sentenceAudio: string;
	    images: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new FieldsMapping(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.firstField = source["firstField"];
	        this.hanzi = source["hanzi"];
	        this.exampleSentence = source["exampleSentence"];
	        this.englishDefinition = source["englishDefinition"];
	        this.chineseDefinition = source["chineseDefinition"];
	        this.pinyin = source["pinyin"];
	        this.hanziAudio = source["hanziAudio"];
	        this.sentenceAudio = source["sentenceAudio"];
	        this.images = source["images"];
	        this.notes = source["notes"];
	    }
	}
	export class AnkiConfig {
	    ActiveDeck: string;
	    ActiveModel: string;
	    ModelMappings: {[key: string]: FieldsMapping};
	    AddProgramTag: boolean;
	    AddBookTag: boolean;
	    AllowDuplicates: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AnkiConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ActiveDeck = source["ActiveDeck"];
	        this.ActiveModel = source["ActiveModel"];
	        this.ModelMappings = this.convertValues(source["ModelMappings"], FieldsMapping, true);
	        this.AddProgramTag = source["AddProgramTag"];
	        this.AddBookTag = source["AddBookTag"];
	        this.AllowDuplicates = source["AllowDuplicates"];
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
	export class Voice {
	    Locale: string;
	    Voice: string;
	    SpeakingStyle: string;
	    RolePlay: string;
	    Speed: number;
	    Pitch: number;
	
	    static createFrom(source: any = {}) {
	        return new Voice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Locale = source["Locale"];
	        this.Voice = source["Voice"];
	        this.SpeakingStyle = source["SpeakingStyle"];
	        this.RolePlay = source["RolePlay"];
	        this.Speed = source["Speed"];
	        this.Pitch = source["Pitch"];
	    }
	}
	export class AzureConfig {
	    GenerateTermAudio: boolean;
	    GenerateSentenceAudio: boolean;
	    AzureApiKey: string;
	    AzureImageApiKey: string;
	    VoiceList: Voice[];
	
	    static createFrom(source: any = {}) {
	        return new AzureConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GenerateTermAudio = source["GenerateTermAudio"];
	        this.GenerateSentenceAudio = source["GenerateSentenceAudio"];
	        this.AzureApiKey = source["AzureApiKey"];
	        this.AzureImageApiKey = source["AzureImageApiKey"];
	        this.VoiceList = this.convertValues(source["VoiceList"], Voice);
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
	export class BookStats {
	    probablyKnownWords: number;
	    knownCharacters: number;
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
	    totalWords: number;
	    uniqueWords: number;
	    totalCharacters: number;
	    uniqueCharacters: number;
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
	        this.totalWords = source["totalWords"];
	        this.uniqueWords = source["uniqueWords"];
	        this.totalCharacters = source["totalCharacters"];
	        this.uniqueCharacters = source["uniqueCharacters"];
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
	export class BookKnownQuery {
	    day: string;
	    known: number;
	
	    static createFrom(source: any = {}) {
	        return new BookKnownQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.day = source["day"];
	        this.known = source["known"];
	    }
	}
	
	export class CalibreBook {
	    authors: string;
	    cover: string;
	    formats: string[];
	    id: number;
	    title: string;
	    exists: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CalibreBook(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.authors = source["authors"];
	        this.cover = source["cover"];
	        this.formats = source["formats"];
	        this.id = source["id"];
	        this.title = source["title"];
	        this.exists = source["exists"];
	    }
	}
	export class CardCreationConfig {
	    AutoAdvanceSentence: boolean;
	    PopulateEnglish: boolean;
	    PopulateChinese: boolean;
	    AutoAdvanceEnglish: boolean;
	    AutoAdvanceChinese: boolean;
	    AutoAdvanceImage: boolean;
	    AutoAdvanceCard: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CardCreationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.AutoAdvanceSentence = source["AutoAdvanceSentence"];
	        this.PopulateEnglish = source["PopulateEnglish"];
	        this.PopulateChinese = source["PopulateChinese"];
	        this.AutoAdvanceEnglish = source["AutoAdvanceEnglish"];
	        this.AutoAdvanceChinese = source["AutoAdvanceChinese"];
	        this.AutoAdvanceImage = source["AutoAdvanceImage"];
	        this.AutoAdvanceCard = source["AutoAdvanceCard"];
	    }
	}
	export class CoverImage {
	    url?: string;
	    imageWidth: number;
	    imageHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new CoverImage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.imageWidth = source["imageWidth"];
	        this.imageHeight = source["imageHeight"];
	    }
	}
	export class Dict {
	    Path: string;
	    Language: string;
	
	    static createFrom(source: any = {}) {
	        return new Dict(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Path = source["Path"];
	        this.Language = source["Language"];
	    }
	}
	export class DictionaryConfig {
	    Dicts: {[key: string]: Dict};
	    PrimaryDict: string;
	    ShowDefinitions: boolean;
	    ShowPinyin: boolean;
	    EnableChinese: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DictionaryConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Dicts = this.convertValues(source["Dicts"], Dict, true);
	        this.PrimaryDict = source["PrimaryDict"];
	        this.ShowDefinitions = source["ShowDefinitions"];
	        this.ShowPinyin = source["ShowPinyin"];
	        this.EnableChinese = source["EnableChinese"];
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
	    pronunciation?: string;
	
	    static createFrom(source: any = {}) {
	        return new DictionaryDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.definition = source["definition"];
	        this.pronunciation = source["pronunciation"];
	    }
	}
	export class DictionaryInfo {
	    name: string;
	    path: string;
	    type: string;
	    isPrimary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DictionaryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.type = source["type"];
	        this.isPrimary = source["isPrimary"];
	    }
	}
	export class ImageInfo {
	    name?: string;
	    url?: string;
	    imageData?: string;
	    imageWidth: number;
	    imageHeight: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.imageData = source["imageData"];
	        this.imageWidth = source["imageWidth"];
	        this.imageHeight = source["imageHeight"];
	    }
	}
	export class Fields {
	    word: string;
	    sentence?: string;
	    englishDefn?: string;
	    chineseDefn?: string;
	    pinyin?: string;
	    images?: ImageInfo[];
	
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
	        this.images = this.convertValues(source["images"], ImageInfo);
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
	
	
	export class LibraryConfig {
	    OnlyFavorites: boolean;
	    HideRead: boolean;
	    DisplayTable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LibraryConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OnlyFavorites = source["OnlyFavorites"];
	        this.HideRead = source["HideRead"];
	        this.DisplayTable = source["DisplayTable"];
	    }
	}
	export class MetaSettings {
	    EnableExperimental: boolean;
	    Ran: number;
	    Theme: string;
	
	    static createFrom(source: any = {}) {
	        return new MetaSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EnableExperimental = source["EnableExperimental"];
	        this.Ran = source["Ran"];
	        this.Theme = source["Theme"];
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
	    word: string;
	    problems: Problems;
	    notes: string;
	    noteId: number;
	
	    static createFrom(source: any = {}) {
	        return new ProblemCard(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.problems = this.convertValues(source["problems"], Problems);
	        this.notes = source["notes"];
	        this.noteId = source["noteId"];
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
	
	export class Sentence {
	    sentence: string;
	    source?: string;
	
	    static createFrom(source: any = {}) {
	        return new Sentence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sentence = source["sentence"];
	        this.source = source["source"];
	    }
	}
	export class SentenceGenerationConfig {
	    IdealSentenceLength: number;
	    KnownInterval: number;
	
	    static createFrom(source: any = {}) {
	        return new SentenceGenerationConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IdealSentenceLength = source["IdealSentenceLength"];
	        this.KnownInterval = source["KnownInterval"];
	    }
	}
	export class TTSVoice {
	    Name: string;
	    DisplayName: string;
	    LocalName: string;
	    ShortName: string;
	    Gender: string;
	    Locale: string;
	    LocaleName: string;
	    StyleList?: string[];
	    SecondaryLocaleList?: string[];
	    SampleRateHertz: string;
	    VoiceType: string;
	    Status: string;
	    ExtendedPropertyMap?: {[key: string]: string};
	    RolePlayList?: string[];
	    WordsPerMinute: string;
	
	    static createFrom(source: any = {}) {
	        return new TTSVoice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.DisplayName = source["DisplayName"];
	        this.LocalName = source["LocalName"];
	        this.ShortName = source["ShortName"];
	        this.Gender = source["Gender"];
	        this.Locale = source["Locale"];
	        this.LocaleName = source["LocaleName"];
	        this.StyleList = source["StyleList"];
	        this.SecondaryLocaleList = source["SecondaryLocaleList"];
	        this.SampleRateHertz = source["SampleRateHertz"];
	        this.VoiceType = source["VoiceType"];
	        this.Status = source["Status"];
	        this.ExtendedPropertyMap = source["ExtendedPropertyMap"];
	        this.RolePlayList = source["RolePlayList"];
	        this.WordsPerMinute = source["WordsPerMinute"];
	    }
	}
	export class TimeQuery {
	    day: string;
	    known: number;
	    knownCharacters: number;
	
	    static createFrom(source: any = {}) {
	        return new TimeQuery(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.day = source["day"];
	        this.known = source["known"];
	        this.knownCharacters = source["knownCharacters"];
	    }
	}
	export class UnknownWordRow {
	    word: string;
	    pinyin: string;
	    occurance: number;
	    frequency: number;
	    definition: string;
	
	    static createFrom(source: any = {}) {
	        return new UnknownWordRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.pinyin = source["pinyin"];
	        this.occurance = source["occurance"];
	        this.frequency = source["frequency"];
	        this.definition = source["definition"];
	    }
	}
	export class WordListsConfig {
	    WordLists: {[key: string]: string};
	    PrimaryWordList: string;
	
	    static createFrom(source: any = {}) {
	        return new WordListsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WordLists = source["WordLists"];
	        this.PrimaryWordList = source["PrimaryWordList"];
	    }
	}
	export class UserSettings {
	    meta: MetaSettings;
	    CardCreation: CardCreationConfig;
	    AnkiConfig: AnkiConfig;
	    AzureConfig: AzureConfig;
	    Dictionaries: DictionaryConfig;
	    WordLists: WordListsConfig;
	    SentenceGeneration: SentenceGenerationConfig;
	    BookLibrary: LibraryConfig;
	
	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], MetaSettings);
	        this.CardCreation = this.convertValues(source["CardCreation"], CardCreationConfig);
	        this.AnkiConfig = this.convertValues(source["AnkiConfig"], AnkiConfig);
	        this.AzureConfig = this.convertValues(source["AzureConfig"], AzureConfig);
	        this.Dictionaries = this.convertValues(source["Dictionaries"], DictionaryConfig);
	        this.WordLists = this.convertValues(source["WordLists"], WordListsConfig);
	        this.SentenceGeneration = this.convertValues(source["SentenceGeneration"], SentenceGenerationConfig);
	        this.BookLibrary = this.convertValues(source["BookLibrary"], LibraryConfig);
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
	
	
	export class WordGridRow {
	    Word: string;
	    Interval: number;
	    LearnedOn: string;
	    InDict: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WordGridRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Word = source["Word"];
	        this.Interval = source["Interval"];
	        this.LearnedOn = source["LearnedOn"];
	        this.InDict = source["InDict"];
	    }
	}
	export class WordListInfo {
	    name: string;
	    path: string;
	    isPrimary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WordListInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isPrimary = source["isPrimary"];
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

