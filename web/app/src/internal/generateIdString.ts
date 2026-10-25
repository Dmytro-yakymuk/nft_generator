const STRING_LENGTH: number = 10;
const STRING_GENERATION_TYPE: number = 16;
const NUMBER_OF_CHARACTERS_FROM_BEGINNING: number = 2;

/** Random string generation with 10 characters. */
export const generateRandomId = () => Math.random().toString(STRING_GENERATION_TYPE).substr(NUMBER_OF_CHARACTERS_FROM_BEGINNING, STRING_LENGTH);
