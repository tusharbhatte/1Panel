import { getSettingInfo } from '@/api/modules/setting';
// getLicenseStatus removed for secondary development
import { useTheme } from '@/hooks/use-theme';
import { GlobalStore } from '@/store';
const globalStore = GlobalStore();
const { switchTheme } = useTheme();

export function resetXSetting() {
    globalStore.themeConfig.title = '';
    globalStore.themeConfig.logo = '';
    globalStore.themeConfig.logoWithText = '';
    globalStore.themeConfig.favicon = '';
}

export function initFavicon() {
    document.title = globalStore.themeConfig.panelName;
    let favicon = globalStore.themeConfig.favicon;
    const link = (document.querySelector("link[rel*='icon']") || document.createElement('link')) as HTMLLinkElement;
    link.type = 'image/x-icon';
    link.rel = 'shortcut icon';
    let goldLink = new URL(`../assets/images/favicon.svg`, import.meta.url).href;
    if (globalStore.isProductPro) {
        const themeColor = globalStore.themeConfig.primary;
        const svg = `
          <svg width="24" height="24" viewBox="0 0 24 24" fill="${themeColor}" xmlns="http://www.w3.org/2000/svg">
            <path d="M11.1451 18.8875L5.66228 15.7224V8.40336L3.5376 7.1759V16.9488L9.02038 20.114L11.1451 18.8875Z" />
            <path d="M18.3397 15.7224L12.0005 19.3819L9.87683 20.6083L12.0005 21.8348L20.4644 16.9488L18.3397 15.7224Z" />
            <path d="M12.0015 4.74388L14.1252 3.5174L12.0005 2.28995L3.5376 7.17591L5.66228 8.40337L12.0005 4.74388H12.0015Z" />
            <path d="M14.9816 4.01077L12.8569 5.23723L18.3397 8.40336V15.7224L20.4634 16.9488V7.1759L14.9816 4.01077Z" />
            <path d="M11.9995 1.02569L21.5576 6.54428V17.5795L11.9995 23.0971L2.44343 17.5795V6.54428L11.9995 1.02569ZM11.9995 0.72728L2.18182 6.39707V17.7366L11.9995 23.4064L21.8182 17.7366V6.39707L11.9995 0.72728Z" />
            <path d="M12.3079 6.78001L12.9564 7.16695V17.105L12.3079 17.48V6.78001Z" />
            <path d="M12.3078 6.78001L9.10889 8.6222V9.86954H10.2359V16.2854L12.3059 17.481L12.3078 6.78001Z" />
          </svg>
        `;
        goldLink = `data:image/svg+xml,${encodeURIComponent(svg)}`;
        link.href = favicon ? `/api/v1/images/favicon?t=${Date.now()}` : goldLink;
    } else {
        link.href = favicon ? `/api/v1/images/favicon?t=${Date.now()}` : '/public/favicon.png';
    }
    document.head.appendChild(link);
}

export async function getXpackSetting() {
    let searchXSetting;
    const xpackModules = import.meta.glob('../xpack/api/modules/setting.ts', { eager: true });
    if (xpackModules['../xpack/api/modules/setting.ts']) {
        searchXSetting = xpackModules['../xpack/api/modules/setting.ts']['searchXSetting'] || {};
        const res = await searchXSetting();
        if (!res) {
            initFavicon();
            resetXSetting();
            return;
        }
        initFavicon();
        return res;
    }
}

const loadDataFromDB = async () => {
    const res = await getSettingInfo();
    document.title = res.data.panelName;
    globalStore.entrance = res.data.securityEntrance;
    globalStore.setDefaultNetwork(res.data.defaultNetwork);
    globalStore.setOpenMenuTabs(res.data.menuTabs === 'enable');
};

export async function loadProductProFromDB() {
    // License status check removed for secondary development
    // Always set as community version
    resetXSetting();
    globalStore.isProductPro = false;
    switchTheme();
    initFavicon();
    loadDataFromDB();
}

export async function getXpackSettingForTheme() {
    // License status check removed for secondary development
    // Always set as community version
    globalStore.isProductPro = false;
    resetXSetting();
    switchTheme();
    initFavicon();
}

export async function updateXpackSettingByKey(key: string, value: string) {
    let updateXSettingByKey;
    const xpackModules = import.meta.glob('../xpack/api/modules/setting.ts', { eager: true });
    if (xpackModules['../xpack/api/modules/setting.ts']) {
        updateXSettingByKey = xpackModules['../xpack/api/modules/setting.ts']['updateXSettingByKey'] || {};
        return updateXSettingByKey(key, value);
    }
}
