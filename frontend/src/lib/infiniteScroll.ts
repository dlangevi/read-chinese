import { ref } from 'vue';
import {
  GetRows, RowCount,
} from '@wailsjs/backend/wordLists';
import type { IGetRowsParams } from 'ag-grid-community';

export class InfiniteScroll {
  rowCountVal = ref(0);

  get rowCount() : number {
    return this.rowCountVal.value;
  }

  getRows(params : IGetRowsParams) {
    GetRows(params.startRow, params.endRow)
      .then(async (data) => {
        const rowCount = await RowCount();
        this.rowCountVal.value = rowCount;
        console.log('Fetched Data',
          data, rowCount);
        params.successCallback(data, rowCount);
      }).catch((err) => {
        console.log('Error fetching data', err);
        params.failCallback();
      });
  }
}
