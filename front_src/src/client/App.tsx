
import "./Header.css";
import "./Main.css";
import "./Footer.css";
import "./App.css";

import { useState, useEffect, useRef } from "react";
import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";
import { go } from "@codemirror/lang-go";
import { json } from "@codemirror/lang-json";
import { nord } from "@fsegurai/codemirror-theme-nord";

import mermaid from "mermaid";

function App() {
    const codeEditorRef = useRef<HTMLDivElement>(null);
    const codeViewRef = useRef<EditorView>(null);

    const tabNames = useRef<string[]>([]);
    const [activeTab, setActiveTab] = useState<number>(-1);

    const mermaidRef = useRef<HTMLDivElement>(null);
    const [mermaids, setMermaids] = useState<string[]>([]);

    const outputEditorRef = useRef<HTMLDivElement>(null);
    const outputViewRef = useRef<EditorView>(null);

    const fileRef = useRef<File>(null);
    const [_fileContent, setFileContent] = useState<string>("");
    const [_fileName, setFileName] = useState<string>("");

    // Initialize mermaid functionality
    useEffect(() => {
        mermaid.initialize({
            startOnLoad: false, securityLevel: 'strict', markdownAutoWrap: false, theme: "base", themeVariables: {
                primaryColor: '#2e3440', primaryTextColor: '#8fbcbb', primaryBorderColor: '#cbd5e1', lineColor: '#cbd5e1', secondaryColor: '#4c566a', tertiaryColor: '#cbd5e1'
            }
        });
    }, []);

    // Render flowchart for the active tab
    const renderMermaid = (index: number) => {
        if (!mermaidRef.current) return;

        const mermaidId = `mermaid-${index}`;

        mermaidRef.current.innerHTML = `<div id="${mermaidId}" class="mermaid">${mermaids[index] || ""}</div>`;

        try {
            mermaid.init(undefined, `#${mermaidId}`);
        } catch (err) {
            console.error(`Mermaid render error!\n${err}`);
        }
    };

    // Call function to render flowchart by condition
    useEffect(() => {
        if (activeTab >= 0 && mermaids.length > activeTab) {
            renderMermaid(activeTab);
        } else if (mermaids.length === 0 && mermaidRef.current) {
            mermaidRef.current.innerHTML = "";
        }
    }, [activeTab, mermaids]);

    // Create code editor view
    useEffect(() => {
        if (!codeEditorRef.current) return;

        const view = new EditorView({
            doc: "",
            parent: codeEditorRef.current,
            extensions: [
                basicSetup,
                EditorState.readOnly.of(true),
                EditorView.editable.of(false),
                EditorView.contentAttributes.of({ tabindex: "0" }),
                go(),
                nord
            ]
        });

        codeViewRef.current = view;

        return () => {
            view.destroy();
        };
    }, []);

    // Create output editor view
    useEffect(() => {
        if (!outputEditorRef.current) return;

        const view = new EditorView({
            doc: "",
            parent: outputEditorRef.current,
            extensions: [
                basicSetup,
                EditorState.readOnly.of(true),
                EditorView.editable.of(false),
                EditorView.contentAttributes.of({ tabindex: "0" }),
                json(),
                nord
            ]
        });

        outputViewRef.current = view;

        return () => {
            view.destroy();
        };
    }, []);

    // Trigger hidden file input
    const handleFileClick = () => {
        const fileInput = document.getElementById("fileInput") as HTMLInputElement;
        fileInput?.click();
    };

    // Read file content and print
    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];

        if (!file)
            return;

        setFileName(file.name);

        fileRef.current = file

        if (!file.name.toLowerCase().endsWith(".go"))
            return;

        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target?.result as string;

            codeViewRef.current?.dispatch({
                changes: {
                    from: 0,
                    to: codeViewRef.current.state.doc.length,
                    insert: content
                }
            });

            setFileContent(content);
        };
        reader.readAsText(file);
    };

    // Upload file to server
    const handleFileUpload = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        const form = new FormData();
        form.append("file", fileRef.current);

        try {
            const res = await fetch("/upload", {
                method: "POST",
                body: form
            });

            if (!res.ok) {
                const err = await res.text();
                console.error(`Upload failed!\n${err}`);
                alert(`Upload failed!\n${err}`);
                return;
            }

            const data = await res.text();
            console.log(`Upload success!\n${data}`);
            alert(`Upload success!\n${data}`);

        } catch (err) {
            console.error(`Upload error!\n${err}`);
            alert(`Upload error\n${err}`);
        }
    };

    // Run inspector and print mermaids
    const handlePrintMermaids = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        try {
            const res = await fetch("/run", {
                method: "GET"
            });

            if (!res.ok) {
                await res.text();
                console.error(`Inspection failed!\n`);
                alert(`Inspection failed!\n`);
                return;
            }

            const data = await res.json();
            console.log(`Inspection success!\n`);
            alert(`Inspection success!\n`);

            // Print JSON file to debug
            outputViewRef.current?.dispatch({
                changes: {
                    from: 0,
                    to: outputViewRef.current.state.doc.length,
                    insert: JSON.stringify(data, null, 4)
                }
            });

            const output = data.output;

            let outMermaid: string = "";
            let outMermaids: Array<string> = [];

            tabNames.current = [];

            // Convert JSON to mermaidJS
            for (const outFuncsName in output) {
                const outFuncs = output[outFuncsName];

                tabNames.current.push(outFuncsName);

                outMermaid += `flowchart TB\n`;

                if (outFuncs.Nodes) {
                    // Clone and sort outNodes
                    const outNodes = [...outFuncs.Nodes].sort((a, b) => {
                        const ai = Number(a.Id);
                        const bi = Number(b.Id);

                        if (Number.isNaN(ai) || Number.isNaN(bi))
                            return String(a.Id).localeCompare(String(b.Id));

                        return ai - bi;
                    });

                    // Convert outNodes to mermaidJS
                    for (let i = 0; i < outNodes.length; i++) {
                        const outNodeID = outNodes[i].Id;
                        const outNodeCode = outNodes[i].Code;
                        const outSafeCode = outNodeCode.replaceAll("`", "#96;").replaceAll("\"", "#34;");
                        const outNodeType = outNodes[i].Node_type;

                        outMermaid += `    id${outNodeID}`

                        if (outNodeType === "basic") {
                            outMermaid += `[\"\`${outSafeCode}\`\"]`;
                        }
                        else if (outNodeType === "cond") {
                            outMermaid += `{\"\`${outSafeCode}\`\"}`;
                        }

                        outMermaid += `\n`;
                    }
                }

                if (outFuncs.Edges) {
                    // Clone and sort outEdges
                    const outEdges = [...outFuncs.Edges].sort((a, b) => {
                        const ai = Number(a.Id);
                        const bi = Number(b.Id);

                        if (Number.isNaN(ai) || Number.isNaN(bi))
                            return String(a.Id).localeCompare(String(b.Id));

                        return ai - bi;
                    });

                    // Convert outEdges to mermaidJS
                    for (let i = 0; i < outEdges.length; i++) {
                        const outEdgeCond = outEdges[i].Label;
                        const outEdgeFrom = outEdges[i].From_node_loc;
                        const outEdgeDest = outEdges[i].To_node_loc;

                        outMermaid += `    id${outEdgeFrom} `;
                        outMermaid += outEdgeCond !== "" ? `-- ${outEdgeCond} --> ` : `--> `;
                        outMermaid += `id${outEdgeDest}\n`;
                    }
                }

                outMermaids.push(outMermaid);
                outMermaid = "";
            }

            // Set mermaids to globally
            setMermaids(outMermaids);

            // Set default tab
            setActiveTab(outMermaids.length ? 0 : -1);

        } catch (_err) {
            console.error(`Inspection error!\n`);
            alert(`Inspection error!\n`);
        }
    };

    // Run inspector and print
    const handleRunInspection = async () => {
        if (!fileRef.current) {
            alert("No file selected!");
            return;
        }

        try {
            // Upload file to server
            handleFileUpload();

            // Run inspector and print mermaids
            handlePrintMermaids();

            // Run inspector again and print debug
            // handlePrintDebug();
            
        } catch (_err) {
            console.error(`Inspection error!\n`);
            alert(`Inspection error!\n`);
        }
    };

    return (
        <>
            <header>
                <div className="headBoxes" id="headTitleBox">
                    <div id="headTitle">TraceInspector</div>
                </div>
                <div className="headBoxes" id="headButtonBox">
                    <button onClick={handleFileClick} className="headButtons" id="openButton">Open</button>
                    <input id="fileInput" type="file" accept=".go" onChange={handleFileChange} style={{ display: "none" }} />
                    <button className="headButtons" id="runButton" onClick={handleRunInspection}>Run</button>
                </div>
                <br /><br /><br /><br /><hr /><br />
            </header>
            <main>
                <div ref={codeEditorRef} className="mainBoxes" id="codeBox"></div>
                <div className="mainBoxes" id="mermaidBox">
                    {mermaids.length > 1 && (
                        <div className="mermaid-tabs">
                            {mermaids.map((_callback, index) => (
                                <button key={index} className={index === activeTab ? "tab active" : "tab"} onClick={() => setActiveTab(index)}>
                                    {`${tabNames.current[index]}`}
                                </button>
                            ))}
                        </div>
                    )}

                    <div ref={mermaidRef} className="mermaid-container" />
                </div>
                <div ref={outputEditorRef} className="mainBoxes" id="logBox"></div>
            </main>
            <footer>
                <div id="placeholder"></div>
            </footer>
        </>
    )
}

export default App;
