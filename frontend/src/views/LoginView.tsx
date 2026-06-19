import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { useAuth } from '../App'
import { KeyRound, User as UserIcon, Loader2 } from 'lucide-react'

export default function LoginView() {
  const { user, refetchAuth } = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  
  const navigate = useNavigate()

  // 如果已经登录，直接跳转到后台设置
  if (user) {
    navigate('/settings', { replace: true })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      await axios.post('/api/auth/login', { username, password })
      refetchAuth()
      navigate('/settings')
    } catch (err: any) {
      if (err.response?.status === 429) {
        setError('登录尝试过于频繁，请 5 分钟后再试')
      } else {
        setError(err.response?.data?.error || '用户名或密码错误')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-md mx-auto px-5 mt-20 animate-fade-in">
      <div className="border border-stone-200 dark:border-neutral-800 bg-white dark:bg-stone-900 p-8 rounded-2xl shadow-sm space-y-6">
        <div className="text-center space-y-2">
          <h1 className="text-2xl font-bold tracking-tight">登入后台管理</h1>
          <p className="text-sm text-neutral-500">做减法，归于文字</p>
        </div>

        {error && (
          <div className="px-4 py-3 rounded-lg bg-red-50 text-red-600 text-sm dark:bg-red-900/20 dark:text-red-400 border border-red-100 dark:border-red-900/50">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-stone-700 dark:text-neutral-300 block">用户名</label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-neutral-400">
                <UserIcon size={18} />
              </span>
              <input
                required
                type="text"
                className="w-full pl-10 pr-4 py-2 border border-stone-200 dark:border-neutral-700 bg-white dark:bg-neutral-800 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none focus:border-transparent"
                placeholder="admin"
                value={username}
                onChange={e => setUsername(e.target.value)}
              />
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-stone-700 dark:text-neutral-300 block">密码</label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-neutral-400">
                <KeyRound size={18} />
              </span>
              <input
                required
                type="password"
                className="w-full pl-10 pr-4 py-2 border border-stone-200 dark:border-neutral-700 bg-white dark:bg-neutral-800 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:outline-none focus:border-transparent"
                placeholder="••••••"
                value={password}
                onChange={e => setPassword(e.target.value)}
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-stone-800 dark:bg-neutral-100 text-white dark:text-stone-900 py-2.5 rounded-lg text-sm font-medium hover:bg-stone-700 dark:hover:bg-neutral-200 transition-colors flex items-center justify-center gap-2 cursor-pointer border-0"
          >
            {loading ? (
              <>
                <Loader2 className="animate-spin" size={16} />
                登入中...
              </>
            ) : (
              '立即登录'
            )}
          </button>
        </form>
      </div>
    </div>
  )
}
