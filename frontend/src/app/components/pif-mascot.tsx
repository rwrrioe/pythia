import { motion } from 'motion/react';

interface PifMascotProps {
  message?: string;
  variant?: 'happy' | 'encouraging' | 'thinking';
  size?: 'sm' | 'md' | 'lg';
}

export function PifMascot({ message, variant = 'happy', size = 'md' }: PifMascotProps) {
  const sizeClasses = {
    sm: 'w-16 h-16',
    md: 'w-24 h-24',
    lg: 'w-32 h-32',
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="flex flex-col items-center gap-3"
    >
      {/* Simple yellow dog illustration */}
      <motion.div
        animate={{ 
          y: [0, -8, 0],
        }}
        transition={{
          duration: 2,
          repeat: Infinity,
          ease: "easeInOut"
        }}
        className={`${sizeClasses[size]} rounded-full bg-gradient-to-br from-yellow-300 to-yellow-500 
                    flex items-center justify-center shadow-lg relative`}
      >
        {/* Dog face */}
        <div className="relative">
          {/* Eyes */}
          <div className="flex gap-3 mb-2">
            <div className="w-2 h-2 bg-gray-800 rounded-full"></div>
            <div className="w-2 h-2 bg-gray-800 rounded-full"></div>
          </div>
          {/* Nose */}
          <div className="w-1.5 h-1.5 bg-gray-800 rounded-full mx-auto"></div>
        </div>
        
        {/* Ears */}
        <div className="absolute -left-2 top-2 w-4 h-6 bg-yellow-400 rounded-full transform -rotate-12"></div>
        <div className="absolute -right-2 top-2 w-4 h-6 bg-yellow-400 rounded-full transform rotate-12"></div>
      </motion.div>
      
      {/* Message */}
      {message && (
        <motion.div
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.2 }}
          className="bg-card border border-border rounded-lg px-4 py-2 shadow-sm max-w-xs"
        >
          <p className="text-sm text-muted-foreground text-center">{message}</p>
        </motion.div>
      )}
    </motion.div>
  );
}
