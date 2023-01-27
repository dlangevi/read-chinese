
async function invoke(action:string, params:any):Promise<any> {
  const response = await fetch('http://localhost:8765', {
    method: 'Post',
    body: JSON.stringify({
      action,
      version: 6,
      params: {
        ...params,
      },
    }),
  });
  return response.json();
}

async function test(word) {
  console.log("hi")
  const noteID = await invoke('findNotes', { query: `Hanzi:${word}` });
  if (noteID.result.length > 1) {
    console.log(`Too many or few notes match ${word}, ${noteID.result}`);
    return 'error';
  }
  const noteInfo = await invoke('notesInfo', {
    notes: noteID.result,
  });
  console.log(JSON.stringify(noteInfo, null, 2))

}

test("嫉妒")

