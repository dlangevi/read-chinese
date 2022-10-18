/**
 * @param { import("knex").Knex } knex
 * @returns { Promise<void> }
 */
exports.up = function up(knex) {
  return knex.schema
    .createTable('words', (table) => {
      table.string('word');
      table.boolean('has_flash_card');
      table.boolean('has_sentence');
      table.integer('interval');
      table.timestamps(true, true);
      table.primary('word');
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
