/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function up(knex) {
  return knex.schema
    .createTable('words', (table) => {
      table.increments('id');
      table.string('word')
        .primary();
      table.boolean('has_flash_card');
      table.boolean('has_sentence');
      table.integer('interval');
      table.timestamps(true, true);
    });
};

/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.down = function down(knex) {
  return knex.schema
    .dropTable('words');
};
