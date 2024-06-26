import { useEffect, useState } from 'react';
import { useConfig, useLogStore } from '@/store/config';
import { SelectGameFile, ImportReg, ExportReg, StartGame, StartGis } from '@go/main/App';
import { EventsOn } from "@runtime";

import './App.css';

function App() {
    const [opReg, setOpReg] = useState("");

    const [game, regs, flush] = useConfig(state => [state.game, state.regs, state.flush]);
    const { gisPath, setGisPath, startGis, setStartGis } = useConfig(state => {
        const { gisPath, setGisPath, startGis, setStartGis } = state;
        return { gisPath, setGisPath, startGis, setStartGis }
    });

    const setLogInfo = useLogStore(state => state.append)

    const handle = async (type: "select" | "import" | "export") => {
        switch (type) {
            // 选择游戏路径
            case "select":
                await SelectGameFile();
                break;
            // 切换账号
            case "import":
                await ImportReg(opReg);
                break;
            // 导出注册表
            case "export":
                await ExportReg(opReg);
                break;
        }
        await flush();
    }

    const startGame = async () => {
        await StartGame();
        if (startGis && gisPath) {
            // 启动游戏后20秒再启动 GIS
            setTimeout(() => {
                StartGis(gisPath);
                setLogInfo("启动 GIS ing...")
            }, 20000)
        }
    }

    return (
        <div id="app">

            <p>
                <label htmlFor="game-path">选择游戏路径: </label>
                <button onClick={() => handle("select")} id="game-path">{game || "点击选择游戏路径"}</button>
            </p>
            <p title="输入GIS路径,打开时运行游戏后启动GIS">
                <label htmlFor="gis-path" >GIS 路径:</label>
                <input type="text" id="gis-path" value={gisPath} onChange={e => setGisPath(e.target.value)} />
                <br />
                <label>启动游戏后是否运行GIS:</label>
                <input type="checkbox" checked={startGis} onChange={e => setStartGis(e.target.checked)} />
            </p>

            <fieldset>
                <legend>输入或选择要操作的账号</legend>
                <p>
                    <input type="text" placeholder="输入导出的文件名,以b或g开头" value={opReg} onChange={e => setOpReg(e.target.value)} />
                    <br />
                    <span>tips: 导出为新注册表时在上侧输入文件名g开头的是官服,b开头的是b服</span>
                </p>
                {
                    regs.map(v => {
                        return (
                            <div key={v}>
                                <input
                                    type="radio"
                                    id={v}
                                    value={v}
                                    name="reg"
                                    checked={opReg === v}
                                    onChange={e => setOpReg(e.target.value)}
                                />
                                <label htmlFor={v}>{v}</label>
                            </div>
                        )
                    })
                }
            </fieldset>


            <p className='btn-s'>
                <button onClick={() => handle("import")}>切换账号</button>
                <button onClick={() => handle("export")}>导出当前注册表</button>
                <button onClick={startGame}>开始游戏</button>
            </p>
            <Log />

        </div>
    )
}

function Log() {
    const { logs, append, clear } = useLogStore();

    useEffect(() => {
        return EventsOn("log", (...infos) => {
            console.log("接收到log消息", infos)
            append(infos);
        });
    }, [append])


    return (
        <section>
            <p className='log-btn'><div>日志区域</div><button onClick={clear}>清除日志</button></p>
            <div className="log">
                <pre>
                    {
                        logs.reverse().join("\n")
                    }
                </pre>

            </div>
        </section>
    )
}

export default App
