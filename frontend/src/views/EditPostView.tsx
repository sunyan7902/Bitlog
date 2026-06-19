import { useState, useEffect } from 'react'
import { useNavigate, useParams, Link } from 'react-router-dom'
import { Save, ArrowLeft, Loader2 } from 'lucide-react'
import { slugify } from 'transliteration'
import axios from 'axios'
import MarkdownEditor from '../components/MarkdownEditor'

interface Post {
  id: string
  slug: string
  title: string
  content: string
  published: boolean
  createdAt: string
  updatedAt: string
}

export default function EditPostView() {
  const { slug } = useParams()
  const [title, setTitle] = useState('')
  const [newSlug, setNewSlug] = useState('')
  const [slugManuallyEdited, setSlugManuallyEdited] = useState(false)
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const [error, setError] = useState('')

  const navigate = useNavigate()

  // 1. 获取要编辑的文章详情
  useEffect(() => {
    axios.get(`/api/posts/${slug}`)
      .then(res => {
        const post: Post = res.data
        setTitle(post.title)
        setNewSlug(post.slug)
        setContent(post.content)
      })
      .catch(err => {
        setError(err.response?.data?.error || '无法获取文章详情')
      })
      .finally(() => {
        setFetching(false)
      })
  }, [slug])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim() || !newSlug.trim()) {
      setError('标题和固定链接不能为空')
      return
    }

    setLoading(true)
    setError('')

    try {
      const contentBase64 = btoa(
        Array.from(new TextEncoder().encode(content), b => String.fromCharCode(b)).join('')
      )
      const res = await axios.put(`/api/posts/${slug}`, {
        title,
        slug: newSlug,
        contentBase64,
      })
      navigate(`/${res.data.slug}`)
    } catch (err: any) {
      setError(err.response?.data?.error || '修改失败')
    } finally {
      setLoading(false)
    }
  }

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const titleVal = e.target.value
    setTitle(titleVal)
    if (!slugManuallyEdited) {
      setNewSlug(slugify(titleVal, { lowercase: true, separator: '-' }))
    }
  }

  if (fetching) {
    return (
      <div className="max-w-3xl mx-auto px-5 py-20 flex items-center justify-center">
        <Loader2 className="animate-spin text-neutral-500" size={24} />
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col min-h-[calc(100dvh-140px)] md:h-[calc(100vh-140px)] h-auto max-w-7xl w-full mx-auto px-5 gap-6 pb-10 md:pb-0">
      <div className="flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate(`/${slug}`)}
            className="p-2 text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100 rounded-md hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors cursor-pointer border-0 bg-transparent"
          >
            <ArrowLeft size={20} />
          </button>
          <h1 className="text-2xl font-bold tracking-tight">编辑博客</h1>
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
          <div className="mt-2">
            <Link to="/blog" className="underline font-medium text-xs">返回博客列表</Link>
          </div>
        </div>
      )}

      <div className="flex gap-4 shrink-0 flex-col md:flex-row">
        <div className="flex-1 flex flex-col gap-2">
          <input
            required
            type="text"
            className="w-full border border-neutral-300 dark:border-neutral-700 bg-transparent p-2.5 text-lg rounded-md focus:ring-2 focus:ring-blue-500 focus:outline-[3px] focus:outline-blue-500/20 cursor-text"
            placeholder="输入文章标题..."
            value={title}
            onChange={handleTitleChange}
          />
        </div>
        <div className="flex-1 flex flex-col gap-2">
          <div className="flex items-center">
            <span className="text-neutral-500 bg-neutral-100 p-2.5 border border-r-0 border-neutral-300 dark:bg-neutral-950 dark:border-neutral-700 rounded-l-md text-sm">/</span>
            <input
              required
              type="text"
              className="border border-neutral-300 dark:border-neutral-700 bg-transparent p-2.5 text-sm rounded-r-md flex-1 focus:ring-2 focus:ring-blue-500 focus:outline-[3px] focus:outline-blue-500/20 cursor-text"
              placeholder="my-first-post"
              value={newSlug}
              onChange={e => {
                setSlugManuallyEdited(true)
                setNewSlug(e.target.value.toLowerCase().replace(/[\s]+/g, '-'))
              }}
            />
          </div>
        </div>
      </div>

      <MarkdownEditor
        content={content}
        onChange={setContent}
      />
    </form>
  )
}
