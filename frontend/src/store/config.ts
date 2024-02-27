import { create } from 'zustand';
import { persist } from 'zustand/middleware';

import { GetGameFile, GetRegs } from '@go/main/App';


type configState = {
    game: string;
    regs: string[];
}

type configAction = {
    flush: () => Promise<void>;
}

export const useConfig = create<configState & configAction>()(
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

            // 从golang 中获取参数
            flush,
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
            set({ logs: [...get().logs, info] });
        }
    }
})