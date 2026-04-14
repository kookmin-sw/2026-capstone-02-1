
import "./Header.css"
import "./Main.css"
import "./Footer.css"
import "./App.css"

import { useState, useEffect, useRef } from 'react';
import { EditorView, basicSetup } from 'codemirror';
import { EditorState } from "@codemirror/state"
import { go } from '@codemirror/lang-go';
import { nord } from '@fsegurai/codemirror-theme-nord'

function App() {
    const editorRef = useRef<HTMLDivElement>(null);
    const viewRef = useRef<EditorView | null>(null);

    useEffect(() => {
        if (!editorRef.current) return;

        // Create the editor view
        const view = new EditorView({
            doc: "",
            parent: editorRef.current,
            extensions: [
                basicSetup,
                EditorState.readOnly.of(true),
                EditorView.editable.of(false),
                EditorView.contentAttributes.of({ tabindex: "0" }),
                go(),
                nord
            ]
        });

        viewRef.current = view;

        // Clean up on unmount
        return () => {
            view.destroy();
        };
    }, []);

    const [fileContent, setFileContent] = useState<string>('');
    const [fileName, setFileName] = useState<string>('');

    const handleFileClick = () => {
        // Trigger the hidden file input
        const fileInput = document.getElementById('fileInput') as HTMLInputElement;
        fileInput?.click();
    };

    const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];

        if (file) {
            setFileName(file.name);

            // Read the file content
            const reader = new FileReader();
            reader.onload = (e) => {
                const content = e.target?.result as string;

                viewRef.current?.dispatch({
                    changes: {
                        from: 0,
                        to: viewRef.current.state.doc.length,
                        insert: content
                    }
                });

                setFileContent(content);
            };
            reader.readAsText(file);
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
                    <input id="fileInput" type="file" accept=".go" onChange={handleFileChange} style={{ display: 'none' }} />
                    <button className="headButtons" id="runButton">Run</button>
                </div>
                <br /><br /><br /><br /><br /><br /><hr /><br />
            </header>
            <main>
                <div ref={editorRef} className="mainBoxes" id="codeBox"></div>
                <div className="mainBoxes" id="graphBox"></div>
                <div className="mainBoxes" id="logBox"></div>
            </main>
            <footer>
                <div id="copyright">&copy; {new Date().getFullYear()} Copyright Reserved</div>
            </footer>
        </>
    )
}

export default App;
