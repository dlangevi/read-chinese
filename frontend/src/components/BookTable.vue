<template>
  <table class="m-auto my-8 table w-4/5 shadow-lg shadow-accent">
    <thead>
      <tr>
        <th>Title</th>
        <th>Author</th>
        <th>Percent Known</th>
        <th>Favorite</th>
        <th>Read</th>
        <th />
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="book in props.books"
        :key="book.bookId"
        class="hover"
        @click="openBook(book)"
      >
        <td>{{ book.title }}</td>
        <td>{{ book.author }}</td>
        <td>
          {{
            (book.stats.totalKnownWords / book.totalWords * 100).toFixed(2)
          }}%
        </td>
        <td @click.stop>
          <input
            v-model="book.favorite"
            type="checkbox"
            class="checkbox-primary checkbox"
            @input="favorite(book.bookId, !book.favorite)"
          >
        </td>
        <td @click.stop>
          <input
            v-model="book.hasRead"
            type="checkbox"
            class="checkbox-primary checkbox"
            @input="markRead(book.bookId, !book.hasRead)"
          >
        </td>
        <td @click.stop>
          <button
            class="btn-error btn-sm btn"
            @click="deleteBook(book.bookId)"
          >
            Delete
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script lang="ts" setup>
import type { backend } from '@wailsjs/models';
import { useRouter } from 'vue-router';
import { DeleteBook, SetFavorite, SetRead } from '@wailsjs/backend/bookLibrary';
const router = useRouter();

function openBook(book : backend.Book) {
  router.push(`/book/${book.bookId}`);
}
function favorite(bookId:number, setTo:boolean) {
  SetFavorite(bookId, setTo);
}
function markRead(bookId:number, setTo:boolean) {
  SetRead(bookId, setTo);
}

function deleteBook(bookId:number) {
  DeleteBook(bookId);
}

const props = defineProps<{
  books: backend.Book[]
}>();

</script>
