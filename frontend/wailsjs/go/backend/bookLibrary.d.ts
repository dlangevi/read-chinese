// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';

export function AddBook(arg1:string,arg2:string,arg3:string,arg4:string):Promise<number>;

export function BookExists(arg1:string,arg2:string):Promise<boolean>;

export function BookPathsPortable():Promise<boolean>;

export function DeleteBook(arg1:number):Promise<void>;

export function ExportDetailedBooks(arg1:string):Promise<void>;

export function FixBookPaths():Promise<Array<backend.Book>>;

export function GetBook(arg1:number):Promise<backend.Book>;

export function GetBookFrequencies(arg1:number):Promise<{[key: string]: number}>;

export function GetBookGraph(arg1:number):Promise<Array<backend.BookKnownQuery>>;

export function GetBooks(arg1:Array<number>):Promise<Array<backend.Book>>;

export function GetDetailedBooks():Promise<Array<backend.Book>>;

export function GetFavoriteFrequencies():Promise<{[key: string]: number}>;

export function HealthCheck():Promise<void>;

export function LearningTarget():Promise<Array<string>>;

export function LearningTargetBook(arg1:number):Promise<Array<string>>;

export function LearningTargetFavorites():Promise<Array<string>>;

export function RecalculateBooks():Promise<void>;

export function SetFavorite(arg1:number,arg2:boolean):Promise<void>;

export function SetRead(arg1:number,arg2:boolean):Promise<void>;

export function TotalRead():Promise<number>;

export function TotalReadChars():Promise<number>;
