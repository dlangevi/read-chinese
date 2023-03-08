// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';

export function AddWord(arg1:string,arg2:number):Promise<Error>;

export function AddWords(arg1:Array<backend.WordEntry>):Promise<Error>;

export function GetOccurances(arg1:Array<string>):Promise<{[key: string]: number}>;

export function GetStatsInfo():Promise<Array<backend.TimeQuery>>;

export function GetUnknownHskWords(arg1:string,arg2:number):Promise<Array<string>>;

export function GetWordStats():Promise<backend.WordStats>;

export function ImportCSVWords(arg1:string):Promise<Error>;
