// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';

export function AddBook(arg1:string,arg2:string,arg3:string,arg4:string):Promise<Error>;

export function BookExists(arg1:string,arg2:string):Promise<boolean>;

export function DeleteBook(arg1:number):Promise<Error>;

export function GetBook(arg1:number):Promise<backend.Book>;

export function GetBooks():Promise<Array<backend.Book>>;

export function GetDetailedBooks(arg1:string):Promise<Array<backend.Book>>;

export function GetSomeBooks(arg1:Array<number>):Promise<Array<backend.Book>>;

export function LearningTarget():Promise<Array<backend.WordOccuranceRow>>;

export function LearningTargetBook(arg1:number):Promise<Array<backend.WordOccuranceRow>>;

export function SetFavorite(arg1:number,arg2:boolean):Promise<Error>;

export function SetRead(arg1:number,arg2:boolean):Promise<Error>;

export function TopUnknownWords(arg1:number,arg2:number):Promise<Array<string>>;

export function TotalRead():Promise<number>;

export function TotalReadChars():Promise<number>;
