const FIRST_ELEMENT_INDEX = 0;
const SECOND_ELEMENT_INDEX = 1;

export const changeFileName = (file: File, name: string) => {
    const blob = file.slice(FIRST_ELEMENT_INDEX, file.size, file.type);

    return new File([blob], `${name}.${file.type.split('/')[SECOND_ELEMENT_INDEX]}`, { type: file.type });
};
