export default function Topbar({ title, onMenuClick }) {
  return (
    <header className="topbar">
      <div className="topbar-left">
        <button type="button" className="hamburger-btn" onClick={onMenuClick} aria-label="Buka menu navigasi">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z" />
          </svg>
        </button>
        <h2 className="topbar-title">{title}</h2>
      </div>
    </header>
  )
}
