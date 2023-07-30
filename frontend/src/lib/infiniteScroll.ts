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
    Promise.all([
      GetRows(params.startRow, params.endRow),
      RowCount(),
    ])
      .then(async ([rows, rowCount]) => {
        this.rowCountVal.value = rowCount;
        params.successCallback(rows, rowCount);
      }).catch((err) => {
        console.error('Error fetching data', err);
        params.failCallback();
      });
  }
}
