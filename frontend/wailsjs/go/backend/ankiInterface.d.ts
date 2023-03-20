// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';

export function ConfigurationCheck():Promise<Error>;

export function CreateAnkiNote(arg1:backend.Fields,arg2:Array<string>):Promise<Error>;

export function GetAnkiNote(arg1:number):Promise<backend.RawAnkiNote>;

export function HealthCheck():Promise<Error>;

export function ImportAnkiKeywords():Promise<Error>;

export function ImportAnkiReviewData():Promise<Error>;

export function LoadDecks():Promise<Array<string>>;

export function LoadModelFields(arg1:string):Promise<Array<string>>;

export function LoadModels():Promise<Array<string>>;

export function LoadProblemCards(arg1:string):Promise<Array<backend.ProblemCard>>;

export function LoadTemplate():Promise<Error>;

export function UpdateNoteFields(arg1:number,arg2:backend.Fields):Promise<Error>;

export function UpdatePinyin(arg1:number):Promise<Error>;

export function UpdateSentenceAudio(arg1:number):Promise<Error>;

export function UpdateWordAudio(arg1:number):Promise<Error>;
