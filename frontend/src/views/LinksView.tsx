import { Link } from 'react-router-dom'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import remarkBreaks from 'remark-breaks'
import { Pencil } from 'lucide-react'
import { useAuth } from '../App'

export default function LinksView() {
  const { user, siteConfig } = useAuth()

  const defaultLinksContent = `No links yet. Add some via the admin panel.`
  const linksContent = siteConfig.friends_content?.trim() || defaultLinksContent

  return (
    <div className="max-w-3xl mx-auto px-5 pt-4 animate-fade-in pb-20 relative">
      {user && (
        <div className="flex justify-end mb-4 absolute top-0 right-5">
          <Link
            to="/links/edit"
            className="inline-flex items-center gap-1.5 text-sm text-neutral-500 hover:text-neutral-900 dark:hover:text-neutral-100 border border-neutral-300 dark:border-neutral-700 px-3 py-1.5 rounded-md hover:bg-neutral-100 dark:hover:bg-neutral-800 transition-colors"
          >
            <Pencil size={14} />
            编辑
          </Link>
        </div>
      )}
      
      <div className={`prose prose-neutral dark:prose-invert min-w-full ${user ? 'mt-12' : ''}`}>
        <ReactMarkdown
          remarkPlugins={[remarkGfm, remarkBreaks]}
          components={{
            table: ({ node, ...props }) => (
              <div className="w-full overflow-x-auto my-6 border border-neutral-200 dark:border-neutral-800 rounded-lg">
                <table {...props} className="min-w-full divide-y divide-neutral-200 dark:divide-neutral-800" />
              </div>
            ),
          }}
        >
          {linksContent}
        </ReactMarkdown>
      </div>
    </div>
  )
}
