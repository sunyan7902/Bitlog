import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Save, ArrowLeft } from 'lucide-react'
import axios from 'axios'
import { useAuth } from '../App'
import MarkdownEditor from '../components/MarkdownEditor'

export default function EditHomeView() {
  const { siteConfig, refetchConfig } = useAuth()
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  
  const navigate = useNavigate()

  useEffect(() => {
    if (siteConfig && siteConfig.home_content) {
      setContent(siteConfig.home_content)
    }
  }, [siteConfig])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      await axios.put('/api/site-config', { home_content: content })
      refetchConfig()
      navigate('/')
    } catch (err: any) {
      setError(err.response?.data?.error || '保存失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col min-h-[calc(100dvh-140px)] md:h-[calc(100vh-140px)] h-auto max-w-7xl w-full mx-auto px-5 gap-6 pb-10 md:pb-0">
      <div className="flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate('/')}
            className="p-2 text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100 rounded-md hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors cursor-pointer border-0 bg-transparent"
          >
            <ArrowLeft size={20} />
          </button>
          <h1 className="text-2xl font-bold tracking-tight">编辑首页内容</h1>
        </div>
        <button
          type="submit"
          disabled={loading}
          className="flex items-center gap-2 bg-neutral-900 text-white px-5 py-2 rounded-md hover:bg-neutral-800 disabled:opacity-50 transition-colors dark:bg-neutral-100 dark:text-neutral-900 dark:hover:bg-neutral-200 cursor-pointer border-0"
        >
          <Save size={16} />
          {loading ? '保存中...' : '保存修改'}
        </button>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 p-3 rounded-md text-sm shrink-0 dark:bg-red-900/20 dark:text-red-400">
          {error}
        </div>
      )}

      <MarkdownEditor
        content={content}
        onChange={setContent}
        placeholder="输入首页展示的 Markdown 文本内容..."
      />
    </form>
  )
}
