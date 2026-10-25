/** Returns current website template Url. */
export const websiteTemplatePath = (websiteId: string | undefined): string => {
    const currentDomain = window.location.origin;

    return `${currentDomain}/api/v0/collections/${websiteId}/mint-site`;
};
