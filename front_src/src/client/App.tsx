
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

function App() {
    const codeEditorRef = useRef<HTMLDivElement>(null);
    const outputEditorRef = useRef<HTMLDivElement>(null);
    const codeViewRef = useRef<EditorView>(null);
    const outputViewRef = useRef<EditorView>(null);

    const fileRef = useRef<File>(null);
    const [fileContent, setFileContent] = useState<string>("");
    const [fileName, setFileName] = useState<string>("");

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
            alert("No file selected");
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
            console.error("Upload error!\n", err);
            alert("Upload error!");
        }
    };

    // Run inspector and print
    const handleRunInspection = async () => {
        if (!fileRef.current) {
            alert("No file selected");
            return;
        }

        try {
            const res = await fetch("/run", {
                method: "GET"
            });

            if (!res.ok) {
                const err = await res.text();
                console.error(`Inspection failed!\n${err}`);
                alert(`Inspection failed!\n${err}`);
                return;
            }

            const data = await res.json();
            console.log(`Inspection success!\n${data.output}`);
            alert(`Inspection success!\n${data.output}`);

            const jsonString = JSON.stringify(data, null, 4);

            outputViewRef.current?.dispatch({
                changes: {
                    from: 0,
                    to: outputViewRef.current.state.doc.length,
                    insert: jsonString
                }
            });

        } catch (err) {
            console.error("Inspection error!\n", err);
            alert("Inspection error!");
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
                    <button className="headButtons" id="uploadButton" onClick={handleFileUpload}>Upload</button>
                    <button className="headButtons" id="runButton" onClick={handleRunInspection}>Run</button>
                </div>
                <br /><br /><br /><br /><br /><br /><hr /><br />
            </header>
            <main>
                <div ref={codeEditorRef} className="mainBoxes" id="codeBox"></div>
                <div className="mainBoxes" id="graphBox"></div>
                <div ref={outputEditorRef} className="mainBoxes" id="logBox"></div>
            </main>
            <footer>
                <div id="copyright">&copy; {new Date().getFullYear()} Copyright Reserved</div>
            </footer>
        </>
    )
}

export default App;
