/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function up(knex) {
  return knex.schema
    .createTable('books', (table) => {
      table.increments('bookId');
      table.string('author');
      table.string('title');
      table.string('cover');
      table.integer('filepath');
      table.primary('bookId');
      table.unique(['author', 'title']);
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function down(knex) {
  return knex.schema.dropTable('books');
};
