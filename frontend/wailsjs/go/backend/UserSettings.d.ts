// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';

export function AddVoice(arg1:backend.Voice):Promise<void>;

export function DeleteDict(arg1:string):Promise<void>;

export function DeleteList(arg1:string):Promise<void>;

export function DeleteMapping(arg1:string):Promise<void>;

export function ExportMapping():Promise<backend.FieldsMapping>;

export function FixSettingsPaths():Promise<void>;

export function GetMapping(arg1:string):Promise<backend.FieldsMapping>;

export function GetTimesRan():Promise<number>;

export function GetUserSettings():Promise<backend.UserSettings>;

export function RemoveVoice(arg1:backend.Voice):Promise<void>;

export function SaveDict(arg1:string,arg2:string,arg3:string):Promise<void>;

export function SaveList(arg1:string,arg2:string):Promise<void>;

export function SetMapping(arg1:string,arg2:backend.FieldsMapping):Promise<void>;

export function SetPrimaryDict(arg1:string):Promise<void>;

export function SetPrimaryList(arg1:string):Promise<void>;

export function SetUserSetting(arg1:string,arg2:string):Promise<void>;

export function SetUserSettingBool(arg1:string,arg2:boolean):Promise<void>;

export function SetUserSettingInt(arg1:string,arg2:number):Promise<void>;

export function SettingsPathsPortable():Promise<boolean>;

export function UpdateTimesRan():Promise<void>;
