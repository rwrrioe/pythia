import React from 'react';

export function GreekPattern({ className = "" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 100 20" xmlns="http://www.w3.org/2000/svg" preserveAspectRatio="none">
      <defs>
        <pattern id="greek-key" x="0" y="0" width="20" height="20" patternUnits="userSpaceOnUse">
          <path d="M0,0 L0,5 L5,5 L5,15 L15,15 L15,5 L20,5 L20,0 L15,0 L15,10 L10,10 L10,0 Z" 
                fill="currentColor" 
                fillOpacity="0.15"/>
        </pattern>
      </defs>
      <rect width="100" height="20" fill="url(#greek-key)" />
    </svg>
  );
}

export function LaurelWreath({ className = "" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
      <path d="M100,20 Q80,40 70,60 Q60,80 55,100 Q60,80 70,60 Q80,40 100,20 M100,20 Q120,40 130,60 Q140,80 145,100 Q140,80 130,60 Q120,40 100,20" 
            fill="none" 
            stroke="currentColor" 
            strokeWidth="2" 
            opacity="0.3"/>
      <ellipse cx="65" cy="50" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(-30 65 50)"/>
      <ellipse cx="55" cy="70" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(-40 55 70)"/>
      <ellipse cx="52" cy="90" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(-50 52 90)"/>
      <ellipse cx="135" cy="50" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(30 135 50)"/>
      <ellipse cx="145" cy="70" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(40 145 70)"/>
      <ellipse cx="148" cy="90" rx="8" ry="12" fill="currentColor" opacity="0.2" transform="rotate(50 148 90)"/>
    </svg>
  );
}

export function Column({ className = "" }: { className?: string }) {
  return (
    <svg className={className} viewBox="0 0 60 200" xmlns="http://www.w3.org/2000/svg">
      <rect x="10" y="0" width="40" height="20" fill="currentColor" opacity="0.1"/>
      <rect x="15" y="20" width="30" height="160" fill="currentColor" opacity="0.05"/>
      <rect x="10" y="180" width="40" height="20" fill="currentColor" opacity="0.1"/>
      <line x1="20" y1="25" x2="20" y2="175" stroke="currentColor" strokeWidth="0.5" opacity="0.1"/>
      <line x1="30" y1="25" x2="30" y2="175" stroke="currentColor" strokeWidth="0.5" opacity="0.1"/>
      <line x1="40" y1="25" x2="40" y2="175" stroke="currentColor" strokeWidth="0.5" opacity="0.1"/>
    </svg>
  );
}
