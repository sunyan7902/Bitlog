import { useParams, Link, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import axios from 'axios'
import ReactMarkdown from 'react-markdown'
import rehypeRaw from 'rehype-raw'
import remarkGfm from 'remark-gfm'
import remarkBreaks from 'remark-breaks'
import { Loader2, Trash2, Pencil } from 'lucide-react'
import { useAuth } from '../App'
import PreBlock from '../components/PreBlock'

interface Post {
  id: string
  slug: string
  title: string
  content: string
  published: boolean
  createdAt: string
  updatedAt: string
}

export default function PostDetailView() {
  const { slug } = useParams()
  const { user } = useAuth()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  // 查询文章详情
  const { data: post, isLoading, error } = useQuery<Post>({
    queryKey: ['post', slug],
    queryFn: () => axios.get(`/api/posts/${slug}`).then(res => res.data),
  })

  // 删除文章
  const deleteMutation = useMutation({
    mutationFn: () => axios.delete(`/api/posts/${slug}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['posts'] })
      navigate('/blog')
    },
  })

  const handleDelete = () => {
    if (window.confirm('确定要删除这篇文章吗？此操作不可撤销！')) {
      deleteMutation.mutate()
    }
  }

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto px-5 py-20 flex items-center justify-center">
        <Loader2 className="animate-spin text-neutral-500" size={24} />
      </div>
    )
  }

  if (error || !post) {
    return (
      <div className="max-w-3xl mx-auto px-5 py-20 text-center space-y-4">
        <h1 className="text-4xl font-extrabold">404</h1>
        <p className="text-neutral-500">抱歉，您查找的文章不存在或已被删除。</p>
        <Link to="/blog" className="text-blue-500 hover:underline">
          返回博客列表
        </Link>
      </div>
    )
  }

  return (
    <article className="space-y-8 animate-fade-in max-w-3xl mx-auto px-5">
      <header className="space-y-4 text-center mt-12 mb-16">
        <h1 className="text-3xl sm:text-4xl font-black tracking-tight leading-tight mx-auto">
          {post.title}
        </h1>
        <div className="text-neutral-500 uppercase tracking-widest text-sm space-x-2">
          <time dateTime={post.createdAt}>
            {new Date(post.createdAt).toLocaleDateString('zh-CN', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
            })}
          </time>
        </div>

        {/* 管理员操作项 */}
        {user && (
          <div className="flex justify-center items-center gap-4 pt-4">
            <Link
              to={`/edit/${post.slug}`}
              className="inline-flex items-center gap-1 text-xs font-semibold uppercase tracking-wider text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100"
            >
              <Pencil size={14} />
              编辑
            </Link>
            <button
              onClick={handleDelete}
              disabled={deleteMutation.isPending}
              className="inline-flex items-center gap-1 text-xs font-semibold uppercase tracking-wider text-neutral-500 hover:text-red-500 cursor-pointer bg-transparent border-0"
            >
              <Trash2 size={14} />
              {deleteMutation.isPending ? '正在删除...' : '删除'}
            </button>
          </div>
        )}
      </header>

      <div className="prose prose-neutral min-w-full dark:prose-invert prose-headings:font-bold prose-a:text-blue-600 dark:prose-a:text-blue-400 prose-img:rounded-xl mx-auto pb-20">
        <ReactMarkdown
          remarkPlugins={[remarkGfm, remarkBreaks]}
          rehypePlugins={[rehypeRaw]}
          components={{
            pre: ({ node, ...props }) => <PreBlock {...props} />,
            table: ({ node, ...props }) => (
              <div className="w-full overflow-x-auto my-6 border border-neutral-200 dark:border-neutral-800 rounded-lg">
                <table {...props} className="min-w-full divide-y divide-neutral-200 dark:divide-neutral-800" />
              </div>
            ),
            iframe: ({ node, ...props }) => {
              let src = props.src || ''
              if (src) {
                try {
                  const url = new URL(src.startsWith('//') ? `https:${src}` : src)
                  if (url.searchParams.has('autoplay')) {
                    url.searchParams.set('autoplay', '0')
                  } else {
                    url.searchParams.append('autoplay', '0')
                  }
                  if (src.includes('bilibili.com')) {
                    url.searchParams.set('high_quality', '1')
                  }
                  src = url.toString()
                } catch {
                  const connector = src.includes('?') ? '&' : '?'
                  if (!src.includes('autoplay=')) {
                    src = `${src}${connector}autoplay=0`
                  } else {
                    src = src.replace(/autoplay=[^&]*/, 'autoplay=0')
                  }
                }
              }
              return (
                <iframe
                  {...props}
                  src={src}
                  className="w-full aspect-video rounded-xl my-8 border border-neutral-200 dark:border-neutral-800 shadow-lg"
                  allow="accelerometer; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowFullScreen
                />
              )
            },
          }}
        >
          {post.content}
        </ReactMarkdown>
      </div>
    </article>
  )
}
