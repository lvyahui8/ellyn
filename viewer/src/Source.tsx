
import  { useState,useEffect } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import { StreamLanguage } from '@codemirror/language';
import { go } from '@codemirror/legacy-modes/mode/go';
import axios from "axios";
import {whiteLight} from '@uiw/codemirror-theme-white'

function SourceView() {
    const [value, setValue] = useState("");
    const [error, setError] = useState(null)

    useEffect(() => {
        axios.get('http://localhost:19898/api/source/0')
            .then(resp => {
                setValue(resp.data)
            })
            .catch(err => {
                setError(err.message)
            })
    },[])

    if (error) {
        return <div>Error: {error}</div>
    }

    return <CodeMirror value={value} height="600px" extensions={[StreamLanguage.define(go)]} theme={whiteLight} editable={false}/>;
}
export default SourceView;