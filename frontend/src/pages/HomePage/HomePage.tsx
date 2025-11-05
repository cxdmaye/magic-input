// filepath: /Volumes/Outlet/code/sample/src/example/magic-input2/App.tsx
import { useState, useEffect } from 'react';
import { Sparkles, History, X, Search } from 'lucide-react';
import { InputInterface, HistoryInputInterface } from '@/core/types';
import { allDetectors } from '@/core/detectors';
import { ResultRenderer } from '@/components/Render';


import { AnalyzeContent, GetAIConfig } from '../../../wailsjs/go/main/App';

// 使用 Wails 生成的类型

const ClipboardAnalyzer = () => {
    const [input, setInput] = useState('');
    const [analysis, setAnalysis] = useState<InputInterface[]>([]);
    const [history, setHistory] = useState<HistoryInputInterface[]>([]);
    const [showHistory, setShowHistory] = useState(false);

    // 从 localStorage 加载历史记录
    useEffect(() => {
        const saved = localStorage.getItem('clipboard-history');
        if (saved) {
            try {
                setHistory(JSON.parse(saved));
            } catch (e) {
                console.error('Failed to load history', e);
            }
        }
    }, []);

    const [isAnalyzing, setIsAnalyzing] = useState(false);
    const [results, setResults] = useState<any[]>([]);
    const [startAiAnalysis, setStartAiAnalysis] = useState(false);

    async function getAiInfo(content: string) {
        try {
            // 检查 AI 配置
            const config = await GetAIConfig()
            if (!config.api_key) {
                alert('请先在设置页面配置 AI API Key')
                return
            }

            const analysisResults = await AnalyzeContent(content)
            setResults(analysisResults || [])
        } catch (err) {
            alert(`分析失败: ${err}`)
        } finally {
            setIsAnalyzing(false)
        }
    }

    // 保存历史记录到 localStorage
    const saveToHistory = (content: string, results: InputInterface[]) => {
        const newItem: HistoryInputInterface = {
            id: Date.now(),
            content: content.slice(0, 200),
            timestamp: new Date().toISOString(),
            results
        };
        const newHistory = [newItem, ...history].slice(0, 50);
        setHistory(newHistory);
        localStorage.setItem('clipboard-history', JSON.stringify(newHistory));
    };

    // 删除历史记录
    const deleteHistoryItem = (id: number) => {
        const newHistory = history.filter(item => item.id !== id);
        setHistory(newHistory);
        localStorage.setItem('clipboard-history', JSON.stringify(newHistory));
    };

    // 清空所有历史记录
    const clearHistory = () => {
        setHistory([]);
        localStorage.removeItem('clipboard-history');
    };

    // 分析内容
    const analyzeContent = (content: string) => {
        if (!content.trim()) {
            setAnalysis([]);
            return;
        }

        const results: InputInterface[] = [];

        // 遍历所有检测器
        for (const detector of allDetectors) {
            const result = detector.test(content);
            if (result) {
                results.push(result);
            }
        }

        // 按置信度排序
        results.sort((a, b) => b.confidence - a.confidence);

        if (results.length === 0 || results[0].type === 'text') {
            setStartAiAnalysis(true);
            getAiInfo(input);
        } else {
            setStartAiAnalysis(false);
        }

        setAnalysis(results);
        saveToHistory(content, results);
    };

    const handleAnalyze = () => {
        analyzeContent(input);

    };

    const loadFromHistory = (item: HistoryInputInterface) => {
        setInput(item.content);
        setAnalysis(item.results);
        setShowHistory(false);
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-purple-50 via-blue-50 to-pink-50">
            {/* 头部 */}
            <header className="bg-white/80 backdrop-blur-sm border-b border-gray-200 sticky top-0 z-10">
                <div className="mx-auto px-6 py-4 flex items-center justify-between pr-[120px]">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-blue-500 flex items-center justify-center shadow-lg">
                            <Sparkles className="w-6 h-6 text-white" />
                        </div>
                        <h1 className="text-xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
                            智能剪贴板分析器
                        </h1>
                    </div>
                    <button
                        onClick={() => setShowHistory(!showHistory)}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer rd-[8px]"
                    >
                        <History className="w-5 h-5" />
                        <span className="text-sm font-medium ">{showHistory ? <>返回查询</> : <>历史记录 ({history.length})</>}</span>
                    </button>
                </div>
            </header>

            <div className="max-w-6xl mx-auto px-6 py-12">
                {!showHistory ? (
                    <div className="space-y-8">
                        {/* 输入区域 */}
                        <div className="bg-white rounded-2xl shadow-lg p-8">
                            <div className="mb-4">
                                <h2 className="text-lg font-semibold text-gray-900 mb-2">粘贴内容进行智能分析</h2>
                                <p className="text-sm text-gray-500">
                                    支持识别 20+ 种格式：图片、网址、JSON、颜色、时间戳、坐标、JWT、UUID 等
                                </p>
                            </div>
                            <textarea
                                value={input}
                                onChange={(e) => setInput(e.target.value)}
                                onKeyDown={(e) => {
                                    if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
                                        handleAnalyze();
                                    }
                                }}
                                placeholder="在此粘贴任意内容... (Ctrl/Cmd + Enter 快速分析)"
                                className="w-full h-40 p-4 box-border border border-gray-200 rounded-xl resize-none focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                            />
                            <div className="mt-4 flex gap-3">
                                <button
                                    onClick={handleAnalyze}
                                    disabled={!input.trim()}
                                    className="flex-1 bg-gradient-to-r from-purple-500 to-blue-500 text-white py-3 rounded-xl font-medium hover:shadow-lg hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center gap-2"
                                >
                                    <Search className="w-5 h-5" />
                                    开始分析
                                </button>
                                <button
                                    onClick={() => {
                                        setInput('');
                                        setAnalysis([]);
                                    }}
                                    disabled={!input && analysis.length === 0}
                                    className="px-6 py-3 border border-gray-300 rounded-xl font-medium hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    清空
                                </button>
                            </div>
                        </div>



                        {/* 分析结果 */}
                        {analysis.length > 0 && (
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h2 className="text-xl font-semibold text-gray-900">
                                        分析结果 <span className="text-sm font-normal text-gray-500">({analysis.length} 个匹配)</span>
                                    </h2>
                                </div>
                                <div className="grid gap-4">
                                    {analysis.map((result, index) => (
                                        <ResultRenderer key={`${result.type}-${index}`} result={result} />
                                    ))}
                                </div>
                            </div>
                        )}

                        {
                            startAiAnalysis && (
                                isAnalyzing ? <>
                                    正在使用 AI 分析...

                                </> : (
                                    <div className="bg-white rounded-2xl shadow-lg p-8">
                                        <div className="flex items-center justify-between mb-6">
                                            <h2 className="text-xl font-semibold text-gray-900 flex items-center gap-2">
                                                🤖 AI 智能分析结果
                                                <span className="text-sm font-normal bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
                                                    {results.length} 个分析
                                                </span>
                                            </h2>
                                        </div>
                                        <div className="space-y-4">
                                            {results.map((result, index) => (
                                                <div key={index} className="relative group">
                                                    <div className="p-6 border border-blue-200 rounded-xl bg-gradient-to-r from-blue-50/80 to-purple-50/80 hover:from-blue-50 hover:to-purple-50 transition-all duration-200 hover:shadow-md">
                                                        {/* 置信度标签 */}
                                                        <div className="absolute top-4 right-4">
                                                            <span className={`px-3 py-1 rounded-full text-sm font-semibold ${parseFloat(result.percent) >= 80
                                                                ? 'bg-green-100 text-green-800 border border-green-200'
                                                                : parseFloat(result.percent) >= 50
                                                                    ? 'bg-yellow-100 text-yellow-800 border border-yellow-200'
                                                                    : 'bg-red-100 text-red-800 border border-red-200'
                                                                }`}>
                                                                {result.percent}
                                                            </span>
                                                        </div>

                                                        {/* 主要内容 */}
                                                        <div className="pr-16">
                                                            {/* 分类标题 */}
                                                            <div className="flex items-center gap-3 mb-3">
                                                                <div className="w-8 h-8 rounded-lg bg-blue-500 flex items-center justify-center text-white text-sm font-bold">
                                                                    {index + 1}
                                                                </div>
                                                                <div>
                                                                    <h3 className="text-lg font-semibold text-gray-900">{result.classify}</h3>
                                                                    <span className="text-sm text-gray-500 font-mono bg-gray-100 px-2 py-0.5 rounded">
                                                                        {result.classification} → {result.type}
                                                                    </span>
                                                                </div>
                                                            </div>

                                                            {/* 分析结果 */}
                                                            <div className="mb-4">
                                                                <div className="flex items-start gap-2 mb-2">
                                                                    <span className="text-blue-600 font-semibold text-sm mt-0.5">🎯 分析结果:</span>
                                                                </div>
                                                                <p className="text-gray-800 leading-relaxed bg-white/70 p-3 rounded-lg border border-blue-100">
                                                                    {result.result}
                                                                </p>
                                                            </div>

                                                            {/* 建议 */}
                                                            <div>
                                                                <div className="flex items-start gap-2 mb-2">
                                                                    <span className="text-purple-600 font-semibold text-sm mt-0.5">💡 建议操作:</span>
                                                                </div>
                                                                <p className="text-gray-700 leading-relaxed bg-purple-50/70 p-3 rounded-lg border border-purple-100">
                                                                    {result.suggestion}
                                                                </p>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                ))
                        }
                    </div>
                ) : (
                    /* 历史记录 */
                    <div className="bg-white rounded-2xl shadow-lg p-8">
                        <div className="flex items-center justify-between mb-6">
                            <h2 className="text-xl font-semibold text-gray-900">
                                历史记录 <span className="text-sm font-normal text-gray-500">({history.length} 条)</span>
                            </h2>
                            {history.length > 0 && (
                                <button
                                    onClick={clearHistory}
                                    className="text-sm text-red-500 hover:text-red-600 font-medium cursor-pointer hover:underline p-x-2 p-y-1 rd-[4px]"
                                >
                                    清空全部
                                </button>
                            )}
                        </div>
                        {history.length === 0 ? (
                            <div className="text-center py-12 text-gray-400">
                                <History className="w-16 h-16 mx-auto mb-4 opacity-50" />
                                <p className="text-lg font-medium">暂无历史记录</p>
                                <p className="text-sm mt-2">开始分析内容后将自动保存历史</p>
                            </div>
                        ) : (
                            <div className="space-y-3">
                                {history.map((item) => (
                                    <div
                                        key={item.id}
                                        className="border border-gray-200  hover:bg-gray-50 hover:border-gray-300 transition-all"
                                    >
                                        <div className="flex items-start justify-between gap-4">
                                            <button
                                                onClick={() => loadFromHistory(item)}
                                                className="flex-1 text-left group rd-[8px] p-4 "
                                            >
                                                <div className="text-gray-700 mb-2 line-clamp-2 group-hover:text-gray-900 rd-[8px] cursor-pointer">
                                                    {item.content}
                                                </div>
                                                <div className="flex items-center gap-2 flex-wrap mb-2">
                                                    {item.results.slice(0, 3).map((r, i) => (
                                                        <span
                                                            key={i}
                                                            className="text-xs bg-purple-100 text-purple-700 px-2 py-1 rounded-full font-medium"
                                                        >
                                                            {r.label}
                                                        </span>
                                                    ))}
                                                    {item.results.length > 3 && (
                                                        <span className="text-xs text-gray-500">
                                                            +{item.results.length - 3} 更多
                                                        </span>
                                                    )}
                                                </div>
                                                <p className="text-xs text-gray-400">
                                                    {new Date(item.timestamp).toLocaleString('zh-CN')}
                                                </p>
                                            </button>
                                            <button
                                                onClick={() => deleteHistoryItem(item.id)}
                                                className="text-gray-400 cursor-pointer hover:text-red-500 transition-colors p-1 rd-[6px]"
                                            >
                                                <X className="w-5 h-5" />
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default ClipboardAnalyzer;