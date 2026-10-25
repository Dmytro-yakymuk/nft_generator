
/** Convert image to base64 string. */
export const convertToBase64 = (file: File): Promise<any> => new Promise((resolve, reject) => {
    const fileReader = new FileReader();

    fileReader.readAsDataURL(file);
    fileReader.onload = () => {
        resolve(fileReader.result);
    };
    fileReader.onerror = (error) => {
        reject(error);
    };
});
