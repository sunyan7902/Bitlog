import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import axios from 'axios'
import { useAuth } from '../App'
import { Search, Loader2 } from 'lucide-react'

interface Post {
  id: string
  slug: string
  title: string
  content: string
  published: boolean
  createdAt: string
  updatedAt: string
}

export default function BlogView() {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  
  const pageParam = searchParams.get('page') || '1'
  const queryParam = searchParams.get('q') || ''
  
  const [searchInput, setSearchInput] = useState(queryParam)

  // 同步搜索输入框的值
  useEffect(() => {
    setSearchInput(queryParam)
  }, [queryParam])

  // React Query 获取博客列表数据
  const { data, isLoading } = useQuery({
    queryKey: ['posts', pageParam, queryParam],
    queryFn: () => axios.get(`/api/posts?page=${pageParam}&q=${queryParam}`).then(res => res.data),
  })

  const posts: Post[] = data?.posts || []
  const total = data?.total || 0
  const currentPage = data?.page || 1
  const totalPages = data?.totalPages || 1

  // 触发搜索
  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setSearchParams({ q: searchInput, page: '1' })
  }

  // 客户端执行年份聚合归档
  const groupedPosts = posts.reduce<Record<string, Post[]>>((acc, post) => {
    const year = new Date(post.createdAt).getFullYear().toString()
    if (!acc[year]) acc[year] = []
    acc[year].push(post)
    return acc
  }, {})

  // 按年份递减排序
  const sortedYears = Object.keys(groupedPosts).sort((a, b) => Number(b) - Number(a))

  if (isLoading) {
    return (
      <div className="max-w-3xl mx-auto px-5 py-20 flex items-center justify-center">
        <Loader2 className="animate-spin text-neutral-500" size={24} />
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto px-5 space-y-10">
      <section className="space-y-8 mt-4">
        {user && (
          <div className="flex flex-col gap-3">
            <form onSubmit={handleSearchSubmit} className="flex items-start justify-end w-full sm:w-80 ml-auto">
              <div className="relative w-full">
                <input
                  type="text"
                  placeholder="搜索博客文章..."
                  className="w-full pl-3 pr-10 py-1.5 border border-neutral-300 dark:border-neutral-700 bg-transparent rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none"
                  value={searchInput}
                  onChange={e => setSearchInput(e.target.value)}
                />
                <button
                  type="submit"
                  className="absolute right-2 top-2 text-neutral-400 hover:text-neutral-900 dark:hover:text-neutral-100 cursor-pointer border-0 bg-transparent"
                >
                  <Search size={16} />
                </button>
              </div>
            </form>
            {queryParam && (
              <p className="text-sm text-neutral-500 text-right">
                搜索 "{queryParam}"，共找到 {total} 篇文章
              </p>
            )}
          </div>
        )}

        {posts.length === 0 ? (
          <p className="text-neutral-500 italic py-8 text-center">目前还没有博客。</p>
        ) : (
          <div className="space-y-8">
            {sortedYears.map(year => (
              <div key={year} className="space-y-1">
                <h2 className="text-3xl font-bold text-stone-300 dark:text-neutral-700 select-none mb-3">
                  {year}
                </h2>
                <div className="space-y-0">
                  {groupedPosts[year].map(post => (
                    <article key={post.id} className="flex items-baseline gap-4 group py-2">
                      <time
                        dateTime={post.createdAt}
                        className="text-neutral-400 dark:text-neutral-500 shrink-0 font-mono text-sm tabular-nums"
                      >
                        {new Date(post.createdAt).toLocaleDateString('zh-CN', {
                          month: 'short',
                          day: '2-digit',
                        })}
                      </time>
                      <Link to={`/${post.slug}`} className="block flex-1 min-w-0">
                        <h3 className="text-base text-blue-600 dark:text-blue-400 group-hover:underline underline-offset-4 decoration-blue-400/50 leading-snug truncate">
                          {post.title}
                        </h3>
                      </Link>
                    </article>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* 分页组件 */}
        {totalPages > 1 && (
          <div className="flex justify-between items-center pt-6 border-t border-neutral-200 dark:border-neutral-800">
            <button
              disabled={currentPage <= 1}
              onClick={() => setSearchParams({ q: queryParam, page: String(currentPage - 1) })}
              className="px-3 py-1.5 border border-neutral-300 dark:border-neutral-700 rounded-md text-sm hover:bg-neutral-100 dark:hover:bg-neutral-800 disabled:opacity-40 transition-colors cursor-pointer bg-transparent"
            >
              上一页
            </button>
            <span className="text-sm text-neutral-500">
              第 {currentPage} / {totalPages} 页
            </span>
            <button
              disabled={currentPage >= totalPages}
              onClick={() => setSearchParams({ q: queryParam, page: String(currentPage + 1) })}
              className="px-3 py-1.5 border border-neutral-300 dark:border-neutral-700 rounded-md text-sm hover:bg-neutral-100 dark:hover:bg-neutral-800 disabled:opacity-40 transition-colors cursor-pointer bg-transparent"
            >
              下一页
            </button>
          </div>
        )}
      </section>
    </div>
  )
}
