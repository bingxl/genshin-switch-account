import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { GetGameFile, GetRegs } from '@go/main/App';


interface configState {
    // 游戏可执行路径
    game: string;
    // 导出的注册表列表
    regs: string[];
    // GIS 可执行路径
    gisPath: string;
    // 启动游戏后是否启动 GIS
    startGis: boolean;


    flush: () => Promise<void>;
    setGisPath: (gisPath: string) => void;
    setStartGis: (startGis: boolean) => void;
}



export const useConfig = create<configState>()(
    persist((set) => {

        const flush = async () => {
            await Promise.all([GetGameFile(), GetRegs()]).then((values) => {
                const [game, regs] = values
                console.log("backend获取到的regs: ", regs)
                set({ game, regs: regs ? [...regs] : [], })
            }).catch(console.error)

        }
        flush();

        return {
            game: '',
            regs: [],
            gisPath: "",
            startGis: false,

            // 从golang 中获取参数
            flush,
            setGisPath: gisPath => set({ gisPath }),
            setStartGis: startGis => set({ startGis }),
        }
    },
        { name: "config-store" }
    )
)


interface LogeI {
    logs: string[];
    append: (...logs: string[]) => void;
    clear: () => void;
}
export const useLogStore = create<LogeI>()((set, get) => {
    return {
        logs: [],
        clear: () => set({ logs: [] }),
        append(...logs) {
            let info = "";
            if (Array.isArray(logs)) {
                info = logs.reduce((pre, cur) => `${pre}\n${cur}`);
            }
            set({ logs: [info, ...get().logs] });
        }
    }
})